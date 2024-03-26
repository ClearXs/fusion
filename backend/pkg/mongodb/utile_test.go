package mongodb

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestValidateLogic(t *testing.T) {
	db := getDb()
	coll := db.Collection("user")
	orPredicate := NewLogicalOr()
	orPredicate.Append(bson.E{Key: "name", Value: "test"})
	orFilter := orPredicate.ToBson()
	_, err := coll.Find(context.TODO(), orFilter)
	assert.Empty(t, err)
}

func TestToJsonString(t *testing.T) {
	a := assert.New(t)

	logical := NewLogical()
	logical.Append(bson.E{Key: "a", Value: "b"})

	r1 := logical.ToJsonString()

	a.Equal("{\"key\":\"$and\",\"value\":{\"a\":\"b\"}}", r1)
}
