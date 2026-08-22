package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// NewsSource — Haber kaynağı
type NewsSource struct {
	Name string `bson:"name" json:"name"`
	URL  string `bson:"url"  json:"url"`
}

// NewsArticle — MongoDB'de cache'lenen haber makalesi
type NewsArticle struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string        `bson:"title"         json:"title"`
	Description string        `bson:"description"   json:"description"`
	Content     string        `bson:"content"       json:"content"`
	URL         string        `bson:"url"           json:"url"`
	Image       string        `bson:"image"         json:"image"`
	PublishedAt string        `bson:"publishedAt"   json:"publishedAt"`
	Source      NewsSource    `bson:"source"        json:"source"`
	Lang        string        `bson:"lang"          json:"lang"`
	CachedAt    time.Time     `bson:"cachedAt"      json:"cachedAt"`
}
