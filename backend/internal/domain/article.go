package domain

import (
	"time"
)

type Article struct {
	Id              uint64    `json:"id" bson:"id"`
	Title           string    `json:"title" bson:"title"`
	Content         string    `json:"content" bson:"content"`
	Tags            []string  `json:"tags" bson:"tags"`
	Top             int64     `json:"top" bson:"top"`
	Category        string    `json:"category" bson:"category"`
	Hidden          bool      `json:"hidden" bson:"hidden"`
	Author          string    `json:"author" bson:"author"`
	Pathname        string    `json:"pathname" bson:"pathname"`
	Private         bool      `json:"private" bson:"private"`
	Password        string    `json:"password" bson:"password"`
	Deleted         bool      `json:"deleted" bson:"deleted"`
	Viewer          int64     `json:"viewer" bson:"viewer"`
	Visited         int64     `json:"visited" bson:"visited"`
	Copyright       string    `json:"copyright" bson:"copyright"`
	LastVisitedTime time.Time `json:"lastVisitedTime" bson:"lastVisitedTime"`
	CreatedAt       time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt" bson:"updatedAt"`
}

type ArticlePageResult struct {
	Articles       []*Article `json:"articles"`
	Total          int64      `json:"total"`
	TotalWordCount int64      `json:"totalWordCount"`
}
