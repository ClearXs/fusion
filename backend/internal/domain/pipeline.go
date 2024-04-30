package domain

import "time"

type EventType = string

const (
	SystemEventType EventType = "system"
	CustomEventType EventType = "custom"
	CornEventType   EventType = "cron"
)

type Pipeline struct {
	Id        uint64    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	EventType EventType `json:"eventType" bson:"eventType"`
	EventName EventKey  `json:"eventName" bson:"eventName"`
	Enabled   bool      `json:"enabled" bson:"enabled"`
	Script    string    `json:"script" bson:"script"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Deleted   bool      `json:"deleted" bson:"deleted"`
}
