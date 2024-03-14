package domain

import (
	"time"
)

type Draft struct {
	Id        int64     `json:"id" bson:"id"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	Tags      []string  `json:"tags" bson:"tags"`
	Category  string    `json:"category" bson:"category"`
	Author    string    `json:"author" bson:"author"`
	Deleted   bool      `json:"deleted" bson:"deleted"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type DraftPageResult struct {
	Drafts []*Draft `json:"articles"`
	Total  int64    `json:"total"`
}
