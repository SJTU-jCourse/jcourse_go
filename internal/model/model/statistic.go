package model

import (
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
	StartDate string `json:"start_time"` // yyyy-mm-dd
	EndDate   string `json:"end_time"`
	Key       string `json:"key"`
	Value     any    `json:"value"`
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
	StartDate      string
	EndDate        string
	PeriodInfoKeys []PeriodInfoKey
}

type UVData = *roaring.Bitmap

type StatisticData struct {
	ID          int64  `json:"id"`
	StatisticID int64  `json:"statistic_id"`
	Date        string `json:"date"` // yyyy-mm-dd
	UVData      UVData `json:"uv_data"`
}
