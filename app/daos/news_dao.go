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

// NewsDAO — MongoDB'de haber cache işlemleri
type NewsDAO struct {
	col *mongo.Collection
}

func NewNewsDAO(db *mongo.Database) *NewsDAO {
	return &NewsDAO{col: db.Collection(collections.NewsArticles)}
}

// FindByLang — Dile göre cache'deki haberleri döner
func (d *NewsDAO) FindByLang(ctx context.Context, lang string) ([]models.NewsArticle, error) {
	filter := bson.M{"lang": lang}
	opts := options.Find().
		SetSort(bson.D{{Key: "cachedAt", Value: -1}}).
		SetLimit(20)

	cursor, err := d.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []models.NewsArticle
	if err := cursor.All(ctx, &articles); err != nil {
		return nil, err
	}

	return articles, nil
}

// BulkUpsert — Eski cache'i temizleyip yeni haberleri toplu olarak kaydeder
func (d *NewsDAO) BulkUpsert(ctx context.Context, articles []models.NewsArticle, lang string) error {
	now := time.Now()

	_, err := d.col.DeleteMany(ctx, bson.M{"lang": lang})
	if err != nil {
		return err
	}

	if len(articles) == 0 {
		return nil
	}

	docs := make([]interface{}, len(articles))
	for i := range articles {
		articles[i].CachedAt = now
		articles[i].Lang = lang
		docs[i] = articles[i]
	}

	_, err = d.col.InsertMany(ctx, docs)
	return err
}
