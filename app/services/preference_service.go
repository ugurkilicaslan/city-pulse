package services

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"city-pulse/app/daos"
	"city-pulse/app/models"
)

// PreferenceService — Kullanıcı tercihleri iş mantığı
type PreferenceService struct {
	dao *daos.PreferenceDAO
}

func NewPreferenceService(dao *daos.PreferenceDAO) *PreferenceService {
	return &PreferenceService{dao: dao}
}

// Get — Kullanıcının tercihlerini getirir, yoksa varsayılan döner
func (s *PreferenceService) Get(ctx context.Context, userID bson.ObjectID) (*models.UserPreference, error) {
	pref, err := s.dao.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if pref == nil {
		return models.DefaultPreference(userID), nil
	}
	return pref, nil
}

// Update — Kullanıcının tercihlerini günceller
func (s *PreferenceService) Update(ctx context.Context, pref *models.UserPreference) error {
	return s.dao.Upsert(ctx, pref)
}

// Reset — Tercihleri varsayılana döndürür
func (s *PreferenceService) Reset(ctx context.Context, userID bson.ObjectID) (*models.UserPreference, error) {
	defaults := models.DefaultPreference(userID)
	return defaults, s.dao.Upsert(ctx, defaults)
}
