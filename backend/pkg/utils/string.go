package utils

import (
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"strconv"
)

const (
	TrueString  = "true"
	FalseString = "false"
)

// ToStringBool 如果字符串是True或者False，否则直接返回false
func ToStringBool(s string) bool {
	if s == TrueString {
		return true
	}
	if s == FalseString {
		return false
	}
	return false
}

// ToStringInt 字符转为int，如果报错则返回 math.MinInt
func ToStringInt(value string, orElse int) int {
	return lo.If[int](value == "", orElse).ElseF(func() int {
		number, err := strconv.Atoi(value)
		if err != nil {
			slog.Error("to string int has err", "err", err.Error())
			return orElse
		}
		return number
	})
}
