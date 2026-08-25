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

// AnalyticsDAO — Kullanıcı davranış olayları için veri erişim katmanı
type AnalyticsDAO struct {
	col *mongo.Collection
}

func NewAnalyticsDAO(db *mongo.Database) *AnalyticsDAO {
	return &AnalyticsDAO{col: db.Collection(collections.AnalyticsEvents)}
}

// Insert — Yeni olay kaydeder
func (d *AnalyticsDAO) Insert(ctx context.Context, e *models.AnalyticsEvent) error {
	e.ID = bson.NewObjectID()
	e.CreatedAt = time.Now()
	_, err := d.col.InsertOne(ctx, e)
	return err
}

// CountByType — Belirli olay tipini sayar
func (d *AnalyticsDAO) CountByType(ctx context.Context, eventType models.EventType, since time.Time) (int64, error) {
	return d.col.CountDocuments(ctx, bson.M{
		"type":       eventType,
		"created_at": bson.M{"$gte": since},
	})
}

// AggregateByType — Olay tipine göre grupla ve say
func (d *AnalyticsDAO) AggregateByType(ctx context.Context, since time.Time) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"created_at": bson.M{"$gte": since}}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$type",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
		{{Key: "$limit", Value: 20}},
	}
	cur, err := d.col.Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var result []bson.M
	return result, cur.All(ctx, &result)
}

// ActiveUsers — Son N dakikada aktif kullanıcı sayısı
func (d *AnalyticsDAO) ActiveUsers(ctx context.Context, withinMinutes int) (int64, error) {
	since := time.Now().Add(-time.Duration(withinMinutes) * time.Minute)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"created_at": bson.M{"$gte": since}}}},
		{{Key: "$group", Value: bson.M{"_id": "$user_id"}}},
		{{Key: "$count", Value: "total"}},
	}
	cur, err := d.col.Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)
	var res []bson.M
	if err := cur.All(ctx, &res); err != nil || len(res) == 0 {
		return 0, err
	}
	if total, ok := res[0]["total"].(int32); ok {
		return int64(total), nil
	}
	return 0, nil
}
