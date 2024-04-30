package util

import (
	"reflect"
)

// Composition generic type T require type must be struct, that don't array or map otherwise ...
// work process from gain outdated same field set it if each latest field is null.
func Composition[T interface{}](outdated *T, latest *T) *T {
	if outdated == nil {
		return latest
	}
	latestValue := reflect.ValueOf(latest)
	outdatedValue := reflect.ValueOf(outdated)
	value := getValue(latestValue)
	for i := 0; i < value.NumField(); i++ {
		lfv := value.Field(i)
		if isNil(lfv) {
			ofv := getValue(outdatedValue).Field(i)
			setValue(lfv, ofv)
		}
	}
	return latest
}

func isNil(value reflect.Value) bool {
	kind := value.Kind()
	switch kind {
	case reflect.Int:
		return value.Int() == 0
	case reflect.Float32:
		return value.Float() == 0
	case reflect.Bool:
		return value.Bool() != true
	case reflect.String:
		return value.String() == ""
	case reflect.Interface:
		return value.Interface() == nil
	case reflect.Map:
		return len(value.MapKeys()) <= 0
	default:
		return false
	}
}

func getValue(origin reflect.Value) reflect.Value {
	if origin.Kind() == reflect.Pointer {
		return origin.Elem()
	}
	return origin
}

func setValue(settable reflect.Value, value reflect.Value) {
	if !settable.CanSet() {
		return
	}
	// compare settable and value kind, if not equal then return and can set
	if settable.Kind() != value.Kind() {
		return
	}
	kind := settable.Kind()
	switch kind {
	case reflect.String:
		settable.SetString(value.String())
	case reflect.Bool:
		settable.SetBool(value.Bool())
	case reflect.Float32:
		settable.SetFloat(value.Float())
	case reflect.Float64:
		settable.SetFloat(value.Float())
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		settable.SetInt(value.Int())
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		settable.SetUint(value.Uint())
	default:
		panic("un discernible type")
	}
}
