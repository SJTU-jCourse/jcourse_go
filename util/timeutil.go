package util

import (
	"time"
)

const (
	GoTimeLayout = "2006-01-02 15:04:05"
	GoDateLayout = "2006-01-02"
)

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(GoTimeLayout, timeStr)
}

func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(GoDateLayout, dateStr)
}

func FormatTime(t time.Time) string {
	return t.Format(GoTimeLayout)
}
func FormatDate(t time.Time) string {
	return t.Format(GoDateLayout)
}
