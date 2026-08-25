package daos

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"city-pulse/app/collections"
	"city-pulse/app/models"
)

// BookmarkDAO — Yer imleri için veri erişim katmanı
type BookmarkDAO struct {
	col *mongo.Collection
}

func NewBookmarkDAO(db *mongo.Database) *BookmarkDAO {
	return &BookmarkDAO{col: db.Collection(collections.Bookmarks)}
}

// FindByUserID — Kullanıcının yer imlerini getirir (son eklenenler önce)
func (d *BookmarkDAO) FindByUserID(ctx context.Context, userID bson.ObjectID) ([]models.Bookmark, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(100)
	cur, err := d.col.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []models.Bookmark
	return result, cur.All(ctx, &result)
}

// FindByType — Belirli tipteki yer imlerini getirir
func (d *BookmarkDAO) FindByType(ctx context.Context, userID bson.ObjectID, btype models.BookmarkType) ([]models.Bookmark, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(50)
	cur, err := d.col.Find(ctx, bson.M{"user_id": userID, "type": btype}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []models.Bookmark
	return result, cur.All(ctx, &result)
}

// Create — Yeni yer imi oluşturur
func (d *BookmarkDAO) Create(ctx context.Context, b *models.Bookmark) error {
	b.ID = bson.NewObjectID()
	b.CreatedAt = time.Now()
	_, err := d.col.InsertOne(ctx, b)
	return err
}

// Delete — Yer imini siler (sadece sahibi silebilir)
func (d *BookmarkDAO) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	res, err := d.col.DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("yer imi bulunamadı")
	}
	return nil
}

// Count — Kullanıcının toplam yer imi sayısı
func (d *BookmarkDAO) Count(ctx context.Context, userID bson.ObjectID) (int64, error) {
	return d.col.CountDocuments(ctx, bson.M{"user_id": userID})
}
