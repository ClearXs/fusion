package logger

// refs:
// https://josephwoodward.co.uk/2022/11/slog-structured-logging-proposal
// https://thedevelopercafe.com/articles/logging-in-go-with-slog-a7bb489755c2

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slog"
)

func GetLogLevel(literal string) slog.Level {
	switch literal {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelWarn
	}
}

type LogrusHandler struct {
	logger *logrus.Logger
}

func NewLogrusHandler(logger *logrus.Logger) *LogrusHandler {
	return &LogrusHandler{logger: logger}
}

func (h *LogrusHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// support all logging levels
	return true
}

func (h *LogrusHandler) Handle(ctx context.Context, rec slog.Record) error {
	fields := make(map[string]interface{}, rec.NumAttrs())

	rec.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	entry := h.logger.WithFields(fields)

	switch rec.Level {
	case slog.LevelDebug:
		entry.Debug(rec.Message)
	case slog.LevelInfo:
		entry.Info(rec.Message)
	case slog.LevelWarn:
		entry.Warn(rec.Message)
	case slog.LevelError:
		entry.Error(rec.Message)
	}

	return nil
}

func (h *LogrusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *LogrusHandler) WithGroup(name string) slog.Handler {
	return h
}
