package util

import (
	"time"
)

const (
	GoTimeLayout = "2006-01-02 15:04:05"
	GoDateLayout = "2006-01-02"
)

func GetLocation() *time.Location {
	name := "Asia/Shanghai" // by default
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

// FormatDate 格式化日期(以环境变量时区为准)
func FormatDate(t time.Time) string {
	dateTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, GetLocation())
	return dateTime.Format(GoDateLayout)
}

// GetDayTimeRange 获取某一天的时间范围(以环境变量时区为准)
func GetDayTimeRange(datetime time.Time) (time.Time, time.Time) {
	year, month, day := datetime.Date()
	location := GetLocation()

	start := time.Date(year, month, day, 0, 0, 0, 0, location)
	end := time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), location)

	return start, end
}

// GetMidTime 获取当天中午12点的时间, 以环境变量指定的时区为准
func GetMidTime(datetime time.Time) time.Time {
	location := GetLocation()
	return time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 12, 0, 0, 0, location)
}
