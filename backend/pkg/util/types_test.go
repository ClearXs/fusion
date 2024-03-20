package util

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestComposition(t *testing.T) {
	d1 := &Demo{Name: "123"}
	d2 := &Demo{Age: 123}
	d2 = Composition[Demo](d1, d2)
	assert.Equal(t, "123", d2.Name)
	assert.Equal(t, 123, d2.Age)

	d3 := &Demo{}
	d3 = Composition[Demo](d1, d3)

	assert.Equal(t, "123", d3.Name)

	d4 := Demo{Name: "123"}
	d5 := Demo{Age: 123}
	d5 = *Composition[Demo](&d4, &d5)

	assert.Equal(t, "123", d5.Name)
	assert.Equal(t, 123, d5.Age)

	d6 := Demo{}
	d6 = *Composition[Demo](&d4, &d6)

	assert.Equal(t, "123", d6.Name)
}
