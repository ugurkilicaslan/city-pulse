package services

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"

	"city-pulse/app/daos"
	"city-pulse/app/models"
)

// ErrAlertNotFound — Uyarı bulunamadı hatası
var ErrAlertNotFound = errors.New("uyarı bulunamadı")

// AlertService — Fiyat uyarıları iş mantığı
type AlertService struct {
	dao *daos.AlertDAO
}

func NewAlertService(dao *daos.AlertDAO) *AlertService {
	return &AlertService{dao: dao}
}

// List — Kullanıcının tüm uyarılarını döner
func (s *AlertService) List(ctx context.Context, userID bson.ObjectID) ([]models.PriceAlert, error) {
	items, err := s.dao.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []models.PriceAlert{}, nil
	}
	return items, nil
}

// Create — Yeni fiyat uyarısı oluşturur
func (s *AlertService) Create(ctx context.Context, a *models.PriceAlert) error {
	return s.dao.Create(ctx, a)
}

// Delete — Fiyat uyarısını siler
func (s *AlertService) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	err := s.dao.Delete(ctx, id, userID)
	if err != nil && err.Error() == "uyarı bulunamadı" {
		return ErrAlertNotFound
	}
	return err
}

// CheckRates — Tüm aktif uyarıları anlık kur ile karşılaştırır
func (s *AlertService) CheckRates(ctx context.Context, rates map[string]float64) ([]models.PriceAlert, error) {
	active, err := s.dao.FindAllActive(ctx)
	if err != nil {
		return nil, err
	}

	var triggered []models.PriceAlert
	for _, alert := range active {
		rate, ok := rates[alert.Pair]
		if !ok {
			continue
		}

		shouldTrigger := false
		switch alert.Condition {
		case models.AlertAbove:
			shouldTrigger = rate >= alert.TargetRate
		case models.AlertBelow:
			shouldTrigger = rate <= alert.TargetRate
		}

		if shouldTrigger {
			_ = s.dao.MarkTriggered(ctx, alert.ID)
			triggered = append(triggered, alert)
		}
	}

	return triggered, nil
}
