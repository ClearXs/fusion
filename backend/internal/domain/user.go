package domain

import (
	"time"
)

type UserType string

const (
	AdminUserType       = "admin"
	CollaborateUserType = "collaborate"
)

type Permission string

type User struct {
	Id          uint64       `json:"id" bson:"id"`
	Name        string       `bson:"name" json:"name"`
	Password    string       `bson:"password" json:"password"`
	CreatedAt   time.Time    `bson:"createdAt" json:"createdAt"`
	Type        UserType     `bson:"type" json:"type"`
	Nickname    string       `bson:"nickname" json:"nickname"`
	Permissions []Permission `bson:"permissions" json:"permissions"`
	Salt        string       `bson:"salt" json:"salt"`
}

type UpdateUser struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}
