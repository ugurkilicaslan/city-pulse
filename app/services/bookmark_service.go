package services

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"

	"city-pulse/app/daos"
	"city-pulse/app/models"
)

// ErrBookmarkNotFound — Yer imi bulunamadı hatası
var ErrBookmarkNotFound = errors.New("yer imi bulunamadı")

// BookmarkService — Yer imleri iş mantığı
type BookmarkService struct {
	dao *daos.BookmarkDAO
}

func NewBookmarkService(dao *daos.BookmarkDAO) *BookmarkService {
	return &BookmarkService{dao: dao}
}

// List — Kullanıcının tüm yer imlerini döner
func (s *BookmarkService) List(ctx context.Context, userID bson.ObjectID) ([]models.Bookmark, error) {
	items, err := s.dao.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []models.Bookmark{}, nil
	}
	return items, nil
}

// ListByType — Belirli tipteki yer imlerini döner
func (s *BookmarkService) ListByType(ctx context.Context, userID bson.ObjectID, btype models.BookmarkType) ([]models.Bookmark, error) {
	items, err := s.dao.FindByType(ctx, userID, btype)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return []models.Bookmark{}, nil
	}
	return items, nil
}

// Add — Yeni yer imi ekler
func (s *BookmarkService) Add(ctx context.Context, b *models.Bookmark) error {
	return s.dao.Create(ctx, b)
}

// Remove — Yer imi siler
func (s *BookmarkService) Remove(ctx context.Context, id, userID bson.ObjectID) error {
	err := s.dao.Delete(ctx, id, userID)
	if err != nil && err.Error() == "yer imi bulunamadı" {
		return ErrBookmarkNotFound
	}
	return err
}

// Count — Toplam yer imi sayısı
func (s *BookmarkService) Count(ctx context.Context, userID bson.ObjectID) (int64, error) {
	return s.dao.Count(ctx, userID)
}
