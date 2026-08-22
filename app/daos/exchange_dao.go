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

// ExchangeDAO — MongoDB'de döviz kuru cache işlemleri
type ExchangeDAO struct {
	col *mongo.Collection
}

func NewExchangeDAO(db *mongo.Database) *ExchangeDAO {
	return &ExchangeDAO{col: db.Collection(collections.ExchangeRates)}
}

// FindByBase — Base currency'ye göre cache'deki son kaydı döner
// Yoksa nil, nil döner (hata değil)
func (d *ExchangeDAO) FindByBase(ctx context.Context, base string) (*models.ExchangeRate, error) {
	filter := bson.M{"base": base}
	opts := options.FindOne().SetSort(bson.D{{Key: "cachedAt", Value: -1}})

	var result models.ExchangeRate
	err := d.col.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

// Upsert — Döviz kurunu MongoDB'ye kaydeder (varsa günceller, yoksa ekler)
func (d *ExchangeDAO) Upsert(ctx context.Context, rate *models.ExchangeRate) error {
	rate.CachedAt = time.Now()

	filter := bson.M{"base": rate.Base}
	update := bson.M{"$set": rate}
	opts := options.UpdateOne().SetUpsert(true)

	_, err := d.col.UpdateOne(ctx, filter, update, opts)
	return err
}
