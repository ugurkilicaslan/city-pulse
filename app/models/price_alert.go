package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// AlertCondition — Fiyat uyarı koşulu
type AlertCondition string

const (
	AlertAbove AlertCondition = "above"
	AlertBelow AlertCondition = "below"
)

// PriceAlert — Kullanıcı fiyat uyarısı
type PriceAlert struct {
	ID          bson.ObjectID  `bson:"_id,omitempty"          json:"id,omitempty"`
	UserID      bson.ObjectID  `bson:"user_id"                json:"userId"`
	Pair        string         `bson:"pair"                   json:"pair"`
	Condition   AlertCondition `bson:"condition"              json:"condition"`
	TargetRate  float64        `bson:"target_rate"            json:"targetRate"`
	Note        string         `bson:"note,omitempty"         json:"note,omitempty"`
	Triggered   bool           `bson:"triggered"              json:"triggered"`
	TriggeredAt *time.Time     `bson:"triggered_at,omitempty" json:"triggeredAt,omitempty"`
	Active      bool           `bson:"active"                 json:"active"`
	CreatedAt   time.Time      `bson:"created_at"             json:"createdAt"`
}
