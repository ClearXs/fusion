package utils

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestWordCount(t *testing.T) {
	count := WordCount("asd撒打算大时，啊，sd")

	assert.Equal(t, 21, count)
}
