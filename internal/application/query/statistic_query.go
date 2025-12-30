package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/infrastructure/entity"
)

type StatisticQueryService interface {
	GetDailyStatistic(ctx context.Context, date string) (*vo.StatisticVO, error)
	GetRangeStatistic(ctx context.Context, startDate string, endDate string) ([]vo.StatisticVO, error)
}

type statisticQueryService struct {
	db *gorm.DB
}

func (s *statisticQueryService) GetDailyStatistic(ctx context.Context, date string) (*vo.StatisticVO, error) {
	st, err := gorm.G[entity.Statistic](s.db).Where("date = ?", date).Take(ctx)
	if err != nil {
		return nil, err
	}

	stVO := vo.NewStatisticVOFromEntity(&st)
	return &stVO, nil
}

func (s *statisticQueryService) GetRangeStatistic(ctx context.Context, startDate string, endDate string) ([]vo.StatisticVO, error) {
	sts, err := gorm.G[entity.Statistic](s.db).Where("date >= ? and date <= ?", startDate, endDate).Find(ctx)
	if err != nil {
		return nil, err
	}

	stVOs := make([]vo.StatisticVO, 0)
	for _, st := range sts {
		stVOs = append(stVOs, vo.NewStatisticVOFromEntity(&st))
	}
	return stVOs, nil
}

func NewStatisticQueryService(db *gorm.DB) StatisticQueryService {
	return &statisticQueryService{db: db}
}
