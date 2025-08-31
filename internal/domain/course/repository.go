package course

import (
	"context"

	"jcourse_go/internal/domain/shared"
)

type CourseRepository interface {
	Get(ctx context.Context, id shared.IDType) (*Course, error)
	GetFilter(ctx context.Context) (*CourseFilter, error)
	FindBy(ctx context.Context, filter CourseListQuery) ([]Course, error)
}

type TeacherRepository interface {
	Get(ctx context.Context, id shared.IDType) (*Teacher, error)
	GetFilter(ctx context.Context) (*TeacherFilter, error)
	FindBy(ctx context.Context, filter TeacherListQuery) ([]Teacher, error)
}

type TrainingPlanRepository interface {
	Get(ctx context.Context, id shared.IDType) (*TrainingPlan, error)
	GetFilter(ctx context.Context) (*TrainingPlanFilter, error)
	FindBy(ctx context.Context, filter TrainingPlanListQuery) ([]TrainingPlan, error)
}

type ReviewRepository interface {
	Get(ctx context.Context, id shared.IDType) (*Review, error)
	GetRevision(ctx context.Context, reviewID shared.IDType) ([]ReviewRevision, error)
	FindBy(ctx context.Context, filter ReviewQuery) ([]Review, error)
	Create(ctx context.Context, review *Review) error
	Update(ctx context.Context, review *Review) error
	Delete(ctx context.Context, review *Review) error
}
