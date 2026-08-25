package ratelimit

import (
	"sync"
	"time"
)

type window struct {
	count   int
	resetAt time.Time
}

// Limiter — Sliding-window rate limiter (IP bazlı)
type Limiter struct {
	mu       sync.Mutex
	visitors map[string]*window
	limit    int
	period   time.Duration
}

// New — limit istek sayısı, period zaman penceresi
func New(limit int, period time.Duration) *Limiter {
	l := &Limiter{
		visitors: make(map[string]*window),
		limit:    limit,
		period:   period,
	}
	go l.cleanup()
	return l
}

// Allow — Verilen key için istek izni verir/vermez
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	w, ok := l.visitors[key]
	if !ok || now.After(w.resetAt) {
		l.visitors[key] = &window{count: 1, resetAt: now.Add(l.period)}
		return true
	}
	if w.count >= l.limit {
		return false
	}
	w.count++
	return true
}

// Remaining — Kalan istek hakkı
func (l *Limiter) Remaining(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	w, ok := l.visitors[key]
	if !ok || time.Now().After(w.resetAt) {
		return l.limit
	}
	rem := l.limit - w.count
	if rem < 0 {
		return 0
	}
	return rem
}

func (l *Limiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		now := time.Now()
		l.mu.Lock()
		for k, w := range l.visitors {
			if now.After(w.resetAt) {
				delete(l.visitors, k)
			}
		}
		l.mu.Unlock()
	}
}
