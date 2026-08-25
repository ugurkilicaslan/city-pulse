package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// UserPreference — Kullanıcıya özel uygulama ayarları
type UserPreference struct {
	ID            bson.ObjectID `bson:"_id,omitempty"   json:"id,omitempty"`
	UserID        bson.ObjectID `bson:"user_id"         json:"userId"`
	DefaultCity   string        `bson:"default_city"    json:"defaultCity"`
	DefaultLang   string        `bson:"default_lang"    json:"defaultLang"`
	CurrencyPairs []string      `bson:"currency_pairs"  json:"currencyPairs"`
	Theme         string        `bson:"theme"           json:"theme"`
	Notifications bool          `bson:"notifications"   json:"notifications"`
	UpdatedAt     time.Time     `bson:"updated_at"      json:"updatedAt"`
}

// DefaultPreference — Yeni kullanıcı için varsayılan tercihler
func DefaultPreference(userID bson.ObjectID) *UserPreference {
	return &UserPreference{
		UserID:        userID,
		DefaultCity:   "istanbul",
		DefaultLang:   "tr",
		CurrencyPairs: []string{"USD/TRY", "EUR/TRY", "CHF/TRY", "GBP/TRY"},
		Theme:         "dark",
		Notifications: false,
		UpdatedAt:     time.Now(),
	}
}
