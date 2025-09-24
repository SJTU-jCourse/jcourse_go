package command

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/statistic"
)

type StatisticService interface {
	CreateTodayStatistic(ctx context.Context) error
}

type statisticService struct {
	db            *gorm.DB
	statisticRepo statistic.StatisticRepository
}

func (s *statisticService) CreateTodayStatistic(ctx context.Context) error {
	return nil
}

func NewStatisticService(
	db *gorm.DB,
	statisticRepo statistic.StatisticRepository,
) StatisticService {
	return &statisticService{
		db:            db,
		statisticRepo: statisticRepo,
	}
}
