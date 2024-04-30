package domain

import "time"

type CustomType = string

const (
	FileCustomType   = "file"
	FolderCustomType = "custom"
)

type CustomPage struct {
	Id        uint64     `json:"id" bson:"id"`
	Name      string     `json:"name" bson:"name"`
	Path      string     `json:"path" bson:"path"`
	Type      CustomType `json:"type" bson:"type"`
	Html      string     `json:"html" bson:"html"`
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt" bson:"updatedAt"`
}
