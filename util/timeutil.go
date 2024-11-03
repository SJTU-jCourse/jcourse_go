package util

import "time"

const layout = "2006-01-02 15:04:05" // Go的参考时间

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}
