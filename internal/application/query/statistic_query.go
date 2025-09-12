package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
)

type StatisticQueryService interface {
	GetDailyStatistic(ctx context.Context, date string) (*vo.StatisticVO, error)
	GetRangeStatistic(ctx context.Context, startDate string, endDate string) ([]vo.StatisticVO, error)
}

type statisticQueryService struct {
	db *gorm.DB
}

func (s *statisticQueryService) GetDailyStatistic(ctx context.Context, date string) (*vo.StatisticVO, error) {
	// TODO implement me
	panic("implement me")
}

func (s *statisticQueryService) GetRangeStatistic(ctx context.Context, startDate string, endDate string) ([]vo.StatisticVO, error) {
	// TODO implement me
	panic("implement me")
}

func NewStatisticQueryService(db *gorm.DB) StatisticQueryService {
	return &statisticQueryService{db: db}
}
