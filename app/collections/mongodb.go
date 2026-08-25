package collections

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoDB collection isimleri — değişmez sabitler
const (
	ExchangeRates   = "exchangeRates"
	NewsArticles    = "newsArticles"
	NasaApods       = "nasaApods"
	Users           = "users"
	AnalyticsEvents = "analyticsEvents"
	UserPreferences = "userPreferences"
	Bookmarks       = "bookmarks"
	PriceAlerts     = "priceAlerts"
)

// Connect — MongoDB'ye bağlanır ve ping atar
func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri).
		SetServerSelectionTimeout(10 * time.Second).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

// EnsureIndexes — TTL index'leri oluşturur (cache için)
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {

	_, err := db.Collection(ExchangeRates).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "cachedAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(3600),
	})
	if err != nil {
		return err
	}

	_, err = db.Collection(NewsArticles).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "cachedAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(900),
	})
	if err != nil {
		return err
	}

	_, err = db.Collection(NasaApods).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "cachedAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(43200),
	})
	if err != nil {
		return err
	}

	_, err = db.Collection(Users).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}
