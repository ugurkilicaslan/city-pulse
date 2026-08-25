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

// AlertDAO — Fiyat uyarıları için veri erişim katmanı
type AlertDAO struct {
	col *mongo.Collection
}

func NewAlertDAO(db *mongo.Database) *AlertDAO {
	return &AlertDAO{col: db.Collection(collections.PriceAlerts)}
}

// FindActiveByUserID — Kullanıcının aktif uyarılarını getirir
func (d *AlertDAO) FindActiveByUserID(ctx context.Context, userID bson.ObjectID) ([]models.PriceAlert, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cur, err := d.col.Find(ctx, bson.M{"user_id": userID, "active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []models.PriceAlert
	return result, cur.All(ctx, &result)
}

// FindAllByUserID — Kullanıcının tüm uyarılarını getirir (tetiklenmiş dahil)
func (d *AlertDAO) FindAllByUserID(ctx context.Context, userID bson.ObjectID) ([]models.PriceAlert, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(50)
	cur, err := d.col.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []models.PriceAlert
	return result, cur.All(ctx, &result)
}

// Create — Yeni fiyat uyarısı oluşturur
func (d *AlertDAO) Create(ctx context.Context, a *models.PriceAlert) error {
	a.ID = bson.NewObjectID()
	a.Active = true
	a.Triggered = false
	a.CreatedAt = time.Now()
	_, err := d.col.InsertOne(ctx, a)
	return err
}

// Delete — Fiyat uyarısını siler (sadece sahibi)
func (d *AlertDAO) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	res, err := d.col.DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("uyarı bulunamadı")
	}
	return nil
}

// MarkTriggered — Uyarıyı tetiklenmiş olarak işaretler ve deaktif eder
func (d *AlertDAO) MarkTriggered(ctx context.Context, id bson.ObjectID) error {
	now := time.Now()
	_, err := d.col.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"triggered":    true,
			"triggered_at": now,
			"active":       false,
		}},
	)
	return err
}

// FindAllActive — Tüm aktif uyarılar (rate checker için)
func (d *AlertDAO) FindAllActive(ctx context.Context) ([]models.PriceAlert, error) {
	opts := options.Find().SetLimit(500)
	cur, err := d.col.Find(ctx, bson.M{"active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []models.PriceAlert
	return result, cur.All(ctx, &result)
}
