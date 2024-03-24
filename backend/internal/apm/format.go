package apm

import (
	"github.com/goccy/go-json"
)

type format = func(data []byte, v interface{}) error

var defaultFormat = func(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Formatter apm info
type Formatter[T interface{}] interface {
	Format() (T, error)
	// FormatTo be without primitive type... it's T must be hava struct type
	FormatTo(v T) error
	FormatToMap() (map[string]any, error)
	FormatToString() (string, error)
}

type SimplyFormatter[T interface{}] struct {
	data []byte
	// format v
	format format
}

func NewFormat[T interface{}](data []byte) *SimplyFormatter[T] {
	return &SimplyFormatter[T]{data: data, format: defaultFormat}
}

func NewFormatByF[T interface{}](data []byte, f format) *SimplyFormatter[T] {
	return &SimplyFormatter[T]{data: data, format: f}
}

func (s *SimplyFormatter[T]) Format() (T, error) {
	t := new(T)
	err := s.format(s.data, t)
	return *t, err
}

func (s *SimplyFormatter[T]) FormatTo(v T) error {
	return s.format(s.data, v)
}

func (s *SimplyFormatter[T]) FormatToMap() (map[string]any, error) {
	var r map[string]any
	err := s.format(s.data, r)
	return r, err
}

func (s *SimplyFormatter[T]) FormatToString() (string, error) {
	var r string
	err := s.format(s.data, r)
	return r, err
}
