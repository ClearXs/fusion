package domain

// SortOrder 排序
type SortOrder = string

const (
	AscSort  SortOrder = "asc"
	DescSort SortOrder = "desc"
)

type Entity struct {
	Id int64 `bson:"id" json:"id"`
}
