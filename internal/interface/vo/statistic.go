package vo

import "jcourse_go/internal/domain/statistic"

type StatisticVO struct {
	Date             string `json:"date"` // yyyy-mm-dd
	NewUserCount     int64  `json:"new_user_count"`
	NewReviewCount   int64  `json:"new_review_count"`
	UVCount          int64  `json:"uv_count"`
	PVCount          int64  `json:"pv_count"`
	TotalUserCount   int64  `json:"total_user_count"`
	TotalReviewCount int64  `json:"total_review_count"`
}

func NewStatisticVO(s *statistic.DailyStatistic) StatisticVO {
	return StatisticVO{
		Date:             s.Date,
		NewUserCount:     s.NewUserCount,
		NewReviewCount:   s.NewReviewCount,
		UVCount:          s.UVCount,
		PVCount:          s.PVCount,
		TotalUserCount:   s.TotalUserCount,
		TotalReviewCount: s.TotalReviewCount,
	}
}
