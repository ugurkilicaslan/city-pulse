package daos

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"city-pulse/app/collections"
	"city-pulse/app/models"
)

// PreferenceDAO — Kullanıcı tercihleri için veri erişim katmanı
type PreferenceDAO struct {
	col *mongo.Collection
}

func NewPreferenceDAO(db *mongo.Database) *PreferenceDAO {
	return &PreferenceDAO{col: db.Collection(collections.UserPreferences)}
}

// FindByUserID — Kullanıcının tercihlerini getirir
func (d *PreferenceDAO) FindByUserID(ctx context.Context, userID bson.ObjectID) (*models.UserPreference, error) {
	var pref models.UserPreference
	err := d.col.FindOne(ctx, bson.M{"user_id": userID}).Decode(&pref)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &pref, err
}

// Upsert — Tercih varsa günceller, yoksa oluşturur
func (d *PreferenceDAO) Upsert(ctx context.Context, pref *models.UserPreference) error {
	pref.UpdatedAt = time.Now()
	opts := options.UpdateOne().SetUpsert(true)
	_, err := d.col.UpdateOne(ctx,
		bson.M{"user_id": pref.UserID},
		bson.M{"$set": pref},
		opts,
	)
	return err
}

// Delete — Kullanıcının tercihlerini siler
func (d *PreferenceDAO) Delete(ctx context.Context, userID bson.ObjectID) error {
	_, err := d.col.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}
