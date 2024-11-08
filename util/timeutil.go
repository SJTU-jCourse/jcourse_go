package util

import (
	"time"
)

const (
	GoTimeLayout = "2006-01-02 15:04:05"
)

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(GoTimeLayout, timeStr)
}
