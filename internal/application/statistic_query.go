package application

import (
	"context"

	"jcourse_go/internal/application/vo"
)

type StatisticQueryService interface {
	GetDailyStatistic(ctx context.Context, date string) (*vo.StatisticVO, error)
	GetRangeStatistic(ctx context.Context, startDate string, endDate string) ([]vo.StatisticVO, error)
}
