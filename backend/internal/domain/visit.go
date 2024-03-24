package domain

import (
	"time"
)

type Visit struct {
	Id              int64     `json:"id" bson:"id"`
	Visited         int64     `json:"visited" bson:"visited"`
	Viewer          int64     `json:"viewer" bson:"viewer"`
	Date            string    `json:"date" bson:"date"`
	Pathname        string    `json:"pathname" bson:"pathname"`
	LastVisitedTime time.Time `json:"lastVisitedTime" bson:"lastVisitedTime"`
	CreatedAt       time.Time `json:"createdAt" bson:"createdAt"`
}
