package domain

type TypeValue[T interface{}] struct {
	Type  string `json:"type" bson:"type"`
	Value T      `json:"value" bson:"value"`
}
