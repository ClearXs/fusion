package mongodb

import (
	"github.com/sony/sonyflake"
	"golang.org/x/exp/slog"
	"strconv"
	"time"
)

var generator *sonyflake.Sonyflake

func init() {
	s := sonyflake.Settings{StartTime: time.Now()}
	generator = sonyflake.NewSonyflake(s)
}

// NextId by snowflake id
func NextId() (uint64, error) {
	id, err := generator.NextID()
	if err != nil {
		slog.Error("Failed to generate snowflake ID by sonyflake", "err", err)
		return 0, err
	}
	return id, nil
}

// NextStringId return string id
func NextStringId() (string, error) {
	id, err := NextId()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), nil
}
