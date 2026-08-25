package metrics

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// EndpointMetrics — Endpoint bazlı metrik kayıtları
type EndpointMetrics struct {
	Calls   int64
	Errors  int64
	TotalMs int64
}

// AvgMs — Ortalama yanıt süresi (ms)
func (e *EndpointMetrics) AvgMs() int64 {
	c := atomic.LoadInt64(&e.Calls)
	if c == 0 {
		return 0
	}
	return atomic.LoadInt64(&e.TotalMs) / c
}

// Registry — Tüm metrikleri tutar
type Registry struct {
	mu            sync.RWMutex
	endpoints     map[string]*EndpointMetrics
	totalRequests int64
	totalErrors   int64
	StartedAt     time.Time
}

// Global — Uygulama genelinde tek kayıt defteri
var Global = &Registry{
	endpoints: make(map[string]*EndpointMetrics),
	StartedAt: time.Now(),
}

// Record — Bir isteği kaydeder
func (r *Registry) Record(endpoint string, latencyMs int64, isError bool) {
	atomic.AddInt64(&r.totalRequests, 1)
	if isError {
		atomic.AddInt64(&r.totalErrors, 1)
	}

	r.mu.Lock()
	em, ok := r.endpoints[endpoint]
	if !ok {
		em = &EndpointMetrics{}
		r.endpoints[endpoint] = em
	}
	r.mu.Unlock()

	atomic.AddInt64(&em.Calls, 1)
	atomic.AddInt64(&em.TotalMs, latencyMs)
	if isError {
		atomic.AddInt64(&em.Errors, 1)
	}
}

// Summary — Tüm metrikleri JSON-serializable olarak döner
func (r *Registry) Summary() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type epSummary struct {
		Endpoint string `json:"endpoint"`
		Calls    int64  `json:"calls"`
		Errors   int64  `json:"errors"`
		AvgMs    int64  `json:"avg_ms"`
	}

	eps := make([]epSummary, 0, len(r.endpoints))
	for path, em := range r.endpoints {
		eps = append(eps, epSummary{
			Endpoint: path,
			Calls:    atomic.LoadInt64(&em.Calls),
			Errors:   atomic.LoadInt64(&em.Errors),
			AvgMs:    em.AvgMs(),
		})
	}
	sort.Slice(eps, func(i, j int) bool { return eps[i].Calls > eps[j].Calls })

	uptimeSec := time.Since(r.StartedAt).Seconds()
	totalReq := atomic.LoadInt64(&r.totalRequests)
	totalErr := atomic.LoadInt64(&r.totalErrors)
	rps := 0.0
	if uptimeSec > 0 {
		rps = float64(totalReq) / uptimeSec
	}

	return map[string]any{
		"uptime_seconds":   uptimeSec,
		"started_at":       r.StartedAt.Format(time.RFC3339),
		"total_requests":   totalReq,
		"total_errors":     totalErr,
		"requests_per_sec": rps,
		"endpoints":        eps,
	}
}
