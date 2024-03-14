package env

import (
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"os"
)

var (
	Version            string
	FusionAllowDomains string
)

func Init() {
	Version = GetEnv("FUSION_VERSION", "")
	FusionAllowDomains = GetEnv("FUSION_ALLOW_DOMAINS", "")
}

// GetEnv get os env
func GetEnv(key string, defaultValue string) string {
	env := os.Getenv(key)
	if lo.IsEmpty(env) {
		return defaultValue
	}
	return env
}

// SetEnv set os env
func SetEnv(key string, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		slog.Error("set env has error", "key", key, "value", value)
	}
}
