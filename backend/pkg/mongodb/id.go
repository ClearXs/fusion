package mongodb

import (
	"github.com/sony/sonyflake"
	"golang.org/x/exp/slog"
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
