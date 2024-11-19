package model

import (
	"time"

	"github.com/RoaringBitmap/roaring"
)

type DailyInfo struct {
	ID               int64  `json:"id"`
	Date             string `json:"date"` // yyyy-mm-dd
	NewUserCount     int64  `json:"new_user_count"`
	NewReviewCount   int64  `json:"new_review_count"`
	UVCount          int64  `json:"uv_count"`
	PVCount          int64  `json:"pv_count"`
	TotalUserCount   int64  `json:"total_user_count"`
	TotalReviewCount int64  `json:"total_review_count"`
}
type PeriodInfo struct {
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Key       string `json:"key"`
	Value     int64  `json:"value"`
}
type PeriodInfoKey = string

const ErrInvalidPeriodInfoKey = "invalid period info key: %s"

const (
	PeriodInfoKeyMAU PeriodInfoKey = "MAU"
	PeriodInfoKeyWAU PeriodInfoKey = "WAU"
)

var PeriodInfoKeys = []PeriodInfoKey{
	PeriodInfoKeyMAU,
	PeriodInfoKeyWAU,
}

type StatisticFilter struct {
	StartTime      time.Time
	EndTime        time.Time
	PeriodInfoKeys []PeriodInfoKey
}

type UVData = *roaring.Bitmap

type StatisticData struct {
	ID          int64  `json:"id"`
	StatisticID int64  `json:"statistic_id"`
	Date        string `json:"date"` // yyyy-mm-dd
	UVData      UVData `json:"uv_data"`
}
