package statistic

import "context"

type StatisticRepository interface {
	Get(ctx context.Context, date string) (*DailyStatistic, error)
	FindByRange(ctx context.Context, startDate string, endDate string) ([]DailyStatistic, error)
	Save(ctx context.Context, stat *DailyStatistic) error
}
