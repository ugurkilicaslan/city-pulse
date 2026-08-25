package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BookmarkType — Kaydedilebilir içerik tipleri
type BookmarkType string

const (
	BookmarkNews BookmarkType = "news"
	BookmarkRepo BookmarkType = "repo"
	BookmarkGame BookmarkType = "game"
	BookmarkNasa BookmarkType = "nasa"
)

// Bookmark — Kullanıcının kaydettiği içerik
type Bookmark struct {
	ID          bson.ObjectID `bson:"_id,omitempty"          json:"id,omitempty"`
	UserID      bson.ObjectID `bson:"user_id"                json:"userId"`
	Type        BookmarkType  `bson:"type"                   json:"type"`
	Title       string        `bson:"title"                  json:"title"`
	URL         string        `bson:"url"                    json:"url"`
	ImageURL    string        `bson:"image_url,omitempty"    json:"imageUrl,omitempty"`
	Description string        `bson:"description,omitempty"  json:"description,omitempty"`
	Source      string        `bson:"source,omitempty"       json:"source,omitempty"`
	Tags        []string      `bson:"tags,omitempty"         json:"tags,omitempty"`
	CreatedAt   time.Time     `bson:"created_at"             json:"createdAt"`
}
