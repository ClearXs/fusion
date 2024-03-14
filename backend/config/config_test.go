package config

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"testing"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type Con struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

func TestSample(test *testing.T) {
	t := Con{}

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t)

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", m)

	d, err = yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}

func TestLocal(t *testing.T) {
	bytes, err := os.ReadFile("./config.yml")

	config := Config{}
	err = yaml.Unmarshal(bytes, &config)

	if err != nil {
		panic(err)
	}
	assert.Equal(t, 5600, config.Server.Port)
}
