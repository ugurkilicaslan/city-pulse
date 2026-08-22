package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ExchangeRate — MongoDB'de saklanan döviz kuru cache modeli
type ExchangeRate struct {
	ID       bson.ObjectID      `bson:"_id,omitempty" json:"id,omitempty"`
	Base     string             `bson:"base"          json:"base"`
	Rates    map[string]float64 `bson:"rates"         json:"rates"`
	Date     string             `bson:"date"          json:"date"`
	CachedAt time.Time          `bson:"cachedAt"      json:"cachedAt"`
}
