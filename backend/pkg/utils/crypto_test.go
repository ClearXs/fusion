package utils

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestEncrypt(t *testing.T) {
	encrypt := Encrypt("username")
	assert.NotEqual(t, "", encrypt)
}

func TestMakeSalt(t *testing.T) {
	salt := MakeSalt()
	assert.NotEqual(t, "", salt)
}
