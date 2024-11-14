package model

import "github.com/RoaringBitmap/roaring"

type DailyInfoMinimal struct {
	ID               int64  `json:"id"`
	Date             string `json:"date"` // yyyy-mm-dd
	NewUserCount     int64  `json:"new_user_count"`
	NewReviewCount   int64  `json:"new_review_count"`
	NewCourseCount   int64  `json:"new_course_count"`
	UVCount          int64  `json:"uv_count"`
	PVCount          int64  `json:"pv_count"`
	TotalUserCount   int64  `json:"total_user_count"`
	TotalReviewCount int64  `json:"total_review_count"`
	TotalCourseCount int64  `json:"total_course_count"`
}

type DailyInfoDetail struct {
	DailyInfoMinimal
	WAU int64 `json:"wau"`
	MAU int64 `json:"mau"`
}

type UVData = *roaring.Bitmap

type StatisticData struct {
	ID          int64  `json:"id"`
	StatisticID int64  `json:"statistic_id"`
	Date        string `json:"date"` // yyyy-mm-dd
	UVData      UVData `json:"uv_data"`
}
