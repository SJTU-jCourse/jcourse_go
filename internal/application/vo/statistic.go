package vo

import (
	"jcourse_go/internal/domain/statistic"
	"jcourse_go/internal/infrastructure/entity"
)

type StatisticVO struct {
	Date             string `json:"date"` // yyyy-mm-dd
	NewUserCount     int64  `json:"new_user_count"`
	NewReviewCount   int64  `json:"new_review_count"`
	UVCount          int64  `json:"uv_count"`
	PVCount          int64  `json:"pv_count"`
	TotalUserCount   int64  `json:"total_user_count"`
	TotalReviewCount int64  `json:"total_review_count"`
}

func NewStatisticVOFromDomain(s *statistic.DailyStatistic) StatisticVO {
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

func NewStatisticVOFromEntity(s *entity.Statistic) StatisticVO {
	return StatisticVO{
		Date:             s.Date,
		NewUserCount:     s.DailyNewUser,
		NewReviewCount:   s.DailyNewReview,
		UVCount:          s.DailyActiveUser,
		PVCount:          s.DailyPageView,
		TotalUserCount:   s.TotalUser,
		TotalReviewCount: s.TotalReview,
	}
}
