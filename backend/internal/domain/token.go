package domain

import (
	"time"
)

type Token struct {
	Id        uint64    `json:"id" bson:"id"`
	UserId    uint64    `json:"userId" bson:"userId"`
	Token     string    `json:"token" bson:"token"`
	Name      string    `json:"name" bson:"name"`
	ExpiresIn int64     `json:"expiresIn" bson:"expiresIn"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Disabled  bool      `json:"disabled" bson:"disabled"`
}
