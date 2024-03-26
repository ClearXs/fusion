package env

import (
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"os"
)

var (
	Version            string
	FusionAllowDomains string
	WebSiteUrl         string
)

func init() {
	Version = GetEnv("FUSION_VERSION", "v0.1.0")
	FusionAllowDomains = GetEnv("FUSION_ALLOW_DOMAINS", "")
	WebSiteUrl = GetEnv("WEB_SITE_URL", "http://127.0.0.1:3001/api/revalidate")
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
