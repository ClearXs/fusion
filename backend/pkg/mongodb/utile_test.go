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
	orFilter := orPredicate.GetBson()
	_, err := coll.Find(context.TODO(), orFilter)
	assert.Empty(t, err)
}
