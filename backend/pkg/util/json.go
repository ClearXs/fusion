package util

import (
	"github.com/goccy/go-json"
	"golang.org/x/exp/slog"
)

// MapToEntity take to map data transfer to the generic type T and return type T
func MapToEntity[T interface{}](data map[string]any, v T) T {
	marshal, err := json.Marshal(data)
	if err != nil {
		slog.Error("map to entity has error", "err", err)
		return v
	}
	if err = json.Unmarshal(marshal, v); err != nil {
		slog.Error("map to entity has error", "err", err)
		return v
	}
	return v
}

// EntityToMap take to generic type T transfer to map
func EntityToMap[T interface{}](v T) map[string]any {
	marshal, err := json.Marshal(v)
	if err != nil {
		slog.Error("entity to map hash error", "err", err)
		return make(map[string]any)
	}
	mapData := make(map[string]any)
	if err = json.Unmarshal(marshal, &mapData); err != nil {
		slog.Error("entity to map hash error", "err", err)
		return mapData
	}
	return mapData
}

// MapArrayToEntityArray take to the array map transfer generic type T for array by MapToEntity
func MapArrayToEntityArray[T interface{}](data []map[string]any, new func() T) []T {
	entities := make([]T, 0)
	for _, v := range data {
		entity := new()
		entity = MapToEntity[T](v, entity)
		entities = append(entities, entity)
	}
	return entities
}

// EntityArrayToMapArray take to generic type T array transfer map array
func EntityArrayToMapArray[T interface{}](values []T) []map[string]any {
	mapData := make([]map[string]any, 0)
	for _, v := range values {
		mapData = append(mapData, EntityToMap[T](v))
	}
	return mapData
}
