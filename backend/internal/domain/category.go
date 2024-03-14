package domain

type CategoryType = string

const (
	CategoryCategoryType CategoryType = "category"
	ColumnCategoryType   CategoryType = "column"
)

type Category struct {
	Id       int64        `json:"id" bson:"id"`
	Name     string       `json:"name" bson:"name"`
	Type     CategoryType `json:"type" bson:"type"`
	Private  bool         `json:"private" bson:"private"`
	Password string       `json:"password" bson:"password"`
	Meta     interface{}  `json:"meta" bson:"meta"`
}
