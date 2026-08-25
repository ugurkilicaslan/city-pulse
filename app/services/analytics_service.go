package services

import (
	"context"
	"time"

	"city-pulse/app/daos"
	"city-pulse/app/models"
	"city-pulse/internal/metrics"
)

// AnalyticsService — Kullanıcı davranışı takip iş mantığı
type AnalyticsService struct {
	dao *daos.AnalyticsDAO
}

func NewAnalyticsService(dao *daos.AnalyticsDAO) *AnalyticsService {
	return &AnalyticsService{dao: dao}
}

// TrackEvent — Olayı asenkron kaydeder (isteği bloklamaz)
func (s *AnalyticsService) TrackEvent(ctx context.Context, e *models.AnalyticsEvent) {
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = s.dao.Insert(bgCtx, e)
	}()
}

// DailySummary — Son 24 saatin olay özetini döner
func (s *AnalyticsService) DailySummary(ctx context.Context) (map[string]any, error) {
	since := time.Now().Add(-24 * time.Hour)

	byType, err := s.dao.AggregateByType(ctx, since)
	if err != nil {
		return nil, err
	}

	activeUsers, err := s.dao.ActiveUsers(ctx, 15)
	if err != nil {
		activeUsers = 0
	}

	total := int64(0)
	for _, row := range byType {
		if c, ok := row["count"].(int32); ok {
			total += int64(c)
		}
	}

	return map[string]any{
		"period":       "last_24h",
		"since":        since.Format(time.RFC3339),
		"total_events": total,
		"active_users": activeUsers,
		"by_type":      byType,
		"server":       metrics.Global.Summary(),
	}, nil
}
