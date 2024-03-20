package util

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

type Demo struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMapToEntity(t *testing.T) {
	data := make(map[string]any)
	data["name"] = "asd"
	data["age"] = 123
	demo := &Demo{}
	demo = MapToEntity[*Demo](data, demo)

	assert.Equal(t, "asd", demo.Name)
	assert.Equal(t, "123", demo.Age)
}

func TestMapToArrayEntity(t *testing.T) {
	data := make(map[string]any)
	data["name"] = "asd"
	data["age"] = 123

	datas := make([]map[string]any, 0)
	datas = append(datas, data)
	demos := MapArrayToEntityArray[*Demo](datas, func() *Demo {
		return &Demo{}
	})
	assert.Equal(t, 1, len(demos))
}

func TestEntityToMap(t *testing.T) {
	demo := &Demo{Name: "asd", Age: 123}
	mapData := EntityToMap[*Demo](demo)
	assert.Equal(t, "asd", mapData["name"])
}

func TestEntityArrayToMap(t *testing.T) {
	demo := &Demo{Name: "asd", Age: 123}
	demos := []*Demo{demo}
	mapData := EntityArrayToMapArray[*Demo](demos)
	assert.Equal(t, 1, len(mapData))
}
