package application

import (
	"context"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/course"
)

type TrainingPlanQueryService interface {
	GetTrainingPlanList(ctx context.Context, query course.TrainingPlanListQuery) ([]vo.TrainingPlanVO, error)
	GetTrainingPlanFilter(ctx context.Context) (*course.TrainingPlanFilter, error)
}
