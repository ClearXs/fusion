package domain

import (
	"time"
)

type StaticType = string
type StorageType = string

const (
	ImgStaticType        StaticType = "img"
	CustomPageStaticType StaticType = "customPage"
)

const (
	PicgoStorageType StorageType = "picgo"
	LocalStorageType StorageType = "local"
)

type Static struct {
	StaticType  StaticType  `json:"staticType" bson:"staticType"`
	StorageType StorageType `json:"storageType" bson:"storageType"`
	FileType    string      `json:"fileType" bson:"fileType"`
	RealPath    string      `json:"realPath" bson:"realPath"`
	Meta        interface{} `json:"meta" bson:"meta"`
	Name        string      `json:"name" bson:"name"`
	Sign        string      `json:"sign" bson:"sign"`
	UpdatedAt   time.Time   `json:"updatedAt" bson:"updatedAt"`
}
