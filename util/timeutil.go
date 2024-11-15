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
func GetDayTimeRange(datetime time.Time) (time.Time, time.Time) {
	year, month, day := datetime.Date()
	location := datetime.Location()

	start := time.Date(year, month, day, 0, 0, 0, 0, location)
	end := time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), location)

	return start, end
}
func GetLocalDayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}
func GetLocalDayEnd() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), time.Local)
}

func GetYesterdayTimeRange(datetime time.Time) (time.Time, time.Time) {
	return GetDayTimeRange(GetTimeSubDay(datetime, 1))
}

func GetTimeSubDay(datetime time.Time, subDay int) time.Time {
	return datetime.AddDate(0, 0, -subDay)
}
func GetTimeSubMonth(datetime time.Time, subMonth int) time.Time {
	return datetime.AddDate(0, -subMonth, 0)
}
func GetTimeSubWeek(datetime time.Time, subWeek int) time.Time {
	return datetime.AddDate(0, 0, -7*subWeek)
}
func GetMidTime(datetime time.Time) time.Time {
	return time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 12, 0, 0, 0, datetime.Location())
}
