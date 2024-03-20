package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Id int `bson:"id"`
}

func TestToBsonMap(t *testing.T) {
	user := User{Id: 1}
	bsonMap := ToBsonMap(&user)
	assert.NotEmpty(t, bsonMap)
	id, ok := bsonMap["id"]
	assert.True(t, ok)
	assert.Equal(t, int32(1), id)
}

func TestToBsonElements(t *testing.T) {
	user := User{Id: 1}
	elements := ToBsonElements(user)
	assert.Equal(t, 1, len(elements))
	e := elements[0]
	assert.Equal(t, "id", e.Key)
	assert.Equal(t, int32(1), e.Value)
}
