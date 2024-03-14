package event

import (
	"github.com/asaskevich/EventBus"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestSubPub(t *testing.T) {
	bus := EventBus.New()

	bus.Subscribe("main:calculator", func(a int, b int) {

		assert.Equal(t, 20, a)
		assert.Equal(t, 40, b)
	})

	bus.Publish("main:calculator", 20, 40)

}
