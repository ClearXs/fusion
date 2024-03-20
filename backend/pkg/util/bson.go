package util

import "go.mongodb.org/mongo-driver/bson"

// ToBsonMap 给定结构体值转换为 bson.M 如果存在错误则panic
func ToBsonMap(val interface{}) bson.M {
	bsonMap := bson.M{}
	bytes, err := bson.Marshal(val)
	if err != nil {
		panic(err)
	}
	if err := bson.Unmarshal(bytes, bsonMap); err != nil {
		panic(err)
	}
	return bsonMap
}

// ToBsonElements 给定的结构转换为 bson.D 如果存在错误则panic
func ToBsonElements(val interface{}) bson.D {
	bsonMap := ToBsonMap(val)
	bsonElements := bson.D{}
	for k, v := range bsonMap {
		bsonElements = append(bsonElements, bson.E{Key: k, Value: v})
	}
	return bsonElements
}
