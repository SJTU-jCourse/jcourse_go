package util

import (
	"time"
)

const (
	GoTimeLayout = "2006-01-02 15:04:05"
	GoDateLayout = "2006-01-02"
)

func GetLocation() *time.Location {
	name := GetTimeLocationStr() // "Asia/Shanghai" by default
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
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

// GetDayTimeRange 获取某一天的时间范围(以环境变量时区为准)
func GetDayTimeRange(datetime time.Time) (time.Time, time.Time) {
	year, month, day := datetime.Date()
	location := GetLocation()

	start := time.Date(year, month, day, 0, 0, 0, 0, location)
	end := time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), location)

	return start, end
}

// GetLocalDayStart 获取当天的开始时间(以Server时区为准)
func GetLocalDayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

// GetLocalDayEnd 获取当天的结束时间(以Server时区为准)
func GetLocalDayEnd() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), time.Local)
}

// GetYesterdayTimeRange 获取昨天的时间范围(以环境变量时区为准)
func GetYesterdayTimeRange(datetime time.Time) (time.Time, time.Time) {
	return GetDayTimeRange(GetTimeSubDay(datetime, 1))
}

// GetTimeSubDay 获取减去指定天数的时间
func GetTimeSubDay(datetime time.Time, subDay int) time.Time {
	return datetime.AddDate(0, 0, -subDay)
}

// GetTimeSubMonth 获取减去指定月数的时间
func GetTimeSubMonth(datetime time.Time, subMonth int) time.Time {
	return datetime.AddDate(0, -subMonth, 0)
}

// GetTimeSubWeek 获取减去指定周数的时间
func GetTimeSubWeek(datetime time.Time, subWeek int) time.Time {
	return datetime.AddDate(0, 0, -7*subWeek)
}

// GetMidTime 获取当天中午12点的时间, 以环境变量指定的时区为准
func GetMidTime(datetime time.Time) time.Time {
	location := GetLocation()
	return time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 12, 0, 0, 0, location)
}
