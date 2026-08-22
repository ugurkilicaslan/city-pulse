package daos

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"city-pulse/app/collections"
	"city-pulse/app/models"
)

type NasaDAO struct {
	collection *mongo.Collection
}

func NewNasaDAO(db *mongo.Database) *NasaDAO {
	return &NasaDAO{
		collection: db.Collection(collections.NasaApods),
	}
}

// GetByDate — Belirli bir güne ait fotoğrafı cache'den getirir
func (d *NasaDAO) GetByDate(ctx context.Context, date string) (*models.ApodEntry, error) {
	var entry models.ApodEntry
	err := d.collection.FindOne(ctx, bson.M{"date": date}).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &entry, nil
}

// Save — NASA verisini cache'e kaydeder
func (d *NasaDAO) Save(ctx context.Context, entry *models.ApodEntry) error {
	opts := options.UpdateOne().SetUpsert(true)
	filter := bson.M{"date": entry.Date}
	update := bson.M{"$set": entry}

	_, err := d.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
