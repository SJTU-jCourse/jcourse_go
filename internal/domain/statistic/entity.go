package statistic

import "time"

type DailyStatistic struct {
	Date             string `json:"date"` // yyyy-mm-dd
	NewUserCount     int64  `json:"new_user_count"`
	NewReviewCount   int64  `json:"new_review_count"`
	UVCount          int64  `json:"uv_count"`
	PVCount          int64  `json:"pv_count"`
	TotalUserCount   int64  `json:"total_user_count"`
	TotalReviewCount int64  `json:"total_review_count"`
	CreatedAt        time.Time
}
