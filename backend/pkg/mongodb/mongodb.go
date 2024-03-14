package mongodb

import (
	"context"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
)

type E struct {
	Key   string
	Value interface{}
}

func Connect(mongodb *Mongodb) (*mongo.Database, func(), error) {
	var uri string
	if lo.IsEmpty(mongodb.Url) {
		uriTemplate := "mongodb://${Username}:${Password}@${Host}:${Port}/?retryWrites=true&w=majority"
		uri = os.Expand(uriTemplate, func(key string) string {
			if key == "Username" {
				return mongodb.Username
			} else if key == "Password" {
				return mongodb.Password
			} else if key == "Host" {
				return mongodb.Host
			} else if key == "Post" {
				return strconv.Itoa(mongodb.Port)
			} else {
				return ""
			}
		})
	} else {
		uri = mongodb.Url
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}
	return client.Database(mongodb.Db), func() { client.Disconnect(context.TODO()) }, nil
}
