package config

import (
	"cc.allio/fusion/internal/apm"
	"cc.allio/fusion/pkg/mongodb"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	Server struct {
		Port int `yaml:"port"`
	}
	Static struct {
		Path string `yaml:"path"`
	}
	Log struct {
		Level string           `yaml:"level"`
		Path  string           `yaml:"path"`
		Apm   apm.LoggerConfig `yaml:"apm"`
	}
	CodeRunner struct {
		Path string `yaml:"path"`
	}
	PluginRunner struct {
		Path string `yaml:"path"`
	}
)

type Config struct {
	Server       Server          `yaml:"server"`
	Static       Static          `yaml:"static"`
	Demo         bool            `yaml:"demo"`
	Log          Log             `yaml:"log"`
	CodeRunner   CodeRunner      `yaml:"codeRunner"`
	PluginRunner PluginRunner    `yaml:"pluginRunner"`
	Mongodb      mongodb.Mongodb `yaml:"mongodb"`
}

func GetConfig() (*Config, error) {
	bytes, err := os.ReadFile("./config/config.yml")
	if err != nil {
		return nil, err
	}
	cfg := Config{}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
