package util

import "golang.org/x/exp/slog"

// GetValue according key and will be force case to generic type, returns defaultValue if error
func GetValue[T interface{}](values map[string]any, key string, defaultValue T) T {
	v, ok := values[key].(T)
	if !ok {
		slog.Debug("Failed to values according to specifies key obtain v", "key", key, "values", values)
		return defaultValue
	}
	return v
}
