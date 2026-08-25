package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// EventType — Takip edilen olay tipleri
type EventType string

const (
	EventPageView    EventType = "page_view"
	EventSectionView EventType = "section_view"
	EventLogin       EventType = "login"
	EventRegister    EventType = "register"
	EventBookmarkAdd EventType = "bookmark_add"
	EventAlertCreate EventType = "alert_create"
	EventSearch      EventType = "search"
	EventRefresh     EventType = "data_refresh"
)

// AnalyticsEvent — Kullanıcı davranış kaydı
type AnalyticsEvent struct {
	ID        bson.ObjectID     `bson:"_id,omitempty"          json:"id,omitempty"`
	UserID    bson.ObjectID     `bson:"user_id"                json:"userId"`
	Type      EventType         `bson:"type"                   json:"type"`
	Page      string            `bson:"page,omitempty"         json:"page,omitempty"`
	Section   string            `bson:"section,omitempty"      json:"section,omitempty"`
	Query     string            `bson:"query,omitempty"        json:"query,omitempty"`
	Metadata  map[string]string `bson:"metadata,omitempty"     json:"metadata,omitempty"`
	IP        string            `bson:"ip,omitempty"           json:"-"`
	UserAgent string            `bson:"user_agent,omitempty"   json:"-"`
	CreatedAt time.Time         `bson:"created_at"             json:"createdAt"`
}
