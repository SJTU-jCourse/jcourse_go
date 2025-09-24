package repository

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/statistic"
	"jcourse_go/internal/infrastructure/entity"
)

type StatisticRepository struct {
	db *gorm.DB
}

func (r *StatisticRepository) Get(ctx context.Context, date string) (*statistic.DailyStatistic, error) {
	e := &entity.Statistic{}
	if err := r.db.Model(&entity.Statistic{}).Where("date = ?", date).First(&e).Error; err != nil {
		return nil, err
	}
	d := newStatisticDomainFromEntity(e)
	return &d, nil
}

func (r *StatisticRepository) FindByRange(ctx context.Context, startDate string, endDate string) ([]statistic.DailyStatistic, error) {
	var entities []entity.Statistic
	if err := r.db.Model(&entity.Statistic{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Order("date ASC").Find(&entities).Error; err != nil {
		return nil, err
	}
	var domains []statistic.DailyStatistic
	for _, e := range entities {
		domains = append(domains, newStatisticDomainFromEntity(&e))
	}
	return domains, nil
}

func (r *StatisticRepository) Save(ctx context.Context, stat *statistic.DailyStatistic) error {
	e := newStatisticEntityFromDomain(stat)
	if err := r.db.Model(&entity.Statistic{}).Where("date = ?", e.Date).
		Create(&e).Error; err != nil {
		return err
	}
	return nil
}

func NewStatisticRepository(db *gorm.DB) statistic.StatisticRepository {
	return &StatisticRepository{db: db}
}

func newStatisticDomainFromEntity(s *entity.Statistic) statistic.DailyStatistic {
	return statistic.DailyStatistic{
		Date:             s.Date,
		NewUserCount:     s.DailyNewUser,
		NewReviewCount:   s.DailyNewReview,
		UVCount:          s.DailyActiveUser,
		PVCount:          s.DailyPageView,
		TotalUserCount:   s.TotalUser,
		TotalReviewCount: s.TotalReview,
		CreatedAt:        s.CreatedAt,
	}
}

func newStatisticEntityFromDomain(s *statistic.DailyStatistic) entity.Statistic {
	return entity.Statistic{
		Date:            s.Date,
		DailyNewUser:    s.NewUserCount,
		DailyNewReview:  s.NewReviewCount,
		DailyActiveUser: s.UVCount,
		DailyPageView:   s.PVCount,
		TotalUser:       s.TotalUserCount,
		TotalReview:     s.TotalReviewCount,
		CreatedAt:       s.CreatedAt,
	}
}
