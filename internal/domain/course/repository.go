package course

import (
	"context"

	"jcourse_go/internal/domain/shared"
)

type CurriculumRepository interface {
	Get(ctx context.Context, code string) (*Curriculum, error)
}

type CourseRepository interface {
	Get(ctx context.Context, id shared.IDType) (*Course, error)
}

type TeacherRepository interface {
	Get(ctx context.Context, id shared.IDType) (*Teacher, error)
}

type TrainingPlanRepository interface {
	Get(ctx context.Context, id shared.IDType) (*TrainingPlan, error)
}

type ReviewRepository interface {
	Get(ctx context.Context, id shared.IDType) (*Review, error)
	Create(ctx context.Context, review *Review) error
	Update(ctx context.Context, review *Review, revision *ReviewRevision) error
	Delete(ctx context.Context, reviewID shared.IDType) error
}
