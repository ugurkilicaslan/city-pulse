package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ApodEntry — NASA Astronomy Picture of the Day, MongoDB cache modeli
type ApodEntry struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string        `bson:"title"         json:"title"`
	Explanation string        `bson:"explanation"   json:"explanation"`
	URL         string        `bson:"url"           json:"url"`
	HdURL       string        `bson:"hdUrl"         json:"hdUrl"`
	MediaType   string        `bson:"mediaType"     json:"mediaType"` // "image" veya "video"
	Date        string        `bson:"date"          json:"date"`
	CachedAt    time.Time     `bson:"cachedAt"      json:"cachedAt"`
}
