package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestFind(t *testing.T) {
	db := getDb()
	coll := db.Collection("user")
	cursor, err := coll.Find(context.TODO(), bson.D{{"name", "test"}})
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			panic(err)
		}
		marshal, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(marshal))
	}
}

func getDb() *mongo.Database {
	mongodb := &Mongodb{Url: "mongodb://root:123456@43.143.195.208:27017/?retryWrites=true&w=majority", Db: "test"}
	db, _, _ := Connect(mongodb)
	return db
}
