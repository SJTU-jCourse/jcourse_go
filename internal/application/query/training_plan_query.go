package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/course"
)

type TrainingPlanQueryService interface {
	GetTrainingPlanList(ctx context.Context, query course.TrainingPlanListQuery) ([]vo.TrainingPlanVO, error)
	GetTrainingPlanFilter(ctx context.Context) (*course.TrainingPlanFilter, error)
}

type trainingPlanQueryService struct {
	db *gorm.DB
}

func (t *trainingPlanQueryService) GetTrainingPlanList(ctx context.Context, query course.TrainingPlanListQuery) ([]vo.TrainingPlanVO, error) {
	// TODO implement me
	panic("implement me")
}

func (t *trainingPlanQueryService) GetTrainingPlanFilter(ctx context.Context) (*course.TrainingPlanFilter, error) {
	// TODO implement me
	panic("implement me")
}

func NewTrainingPlanQueryService(db *gorm.DB) TrainingPlanQueryService {
	return &trainingPlanQueryService{db: db}
}
