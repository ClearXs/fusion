package domain

type File struct {
	Id   int64  `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Path string `json:"path" bson:"path"`
}
