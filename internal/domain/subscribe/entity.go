package subscribe

import (
	"context"
	"time"

	"jcourse_go/internal/domain/shared"
)

type SubscribeLevel string

const (
	SubscribeLevelWatch  SubscribeLevel = "watch"
	SubscribeLevelIgnore SubscribeLevel = "ignore"
)

type CourseSubscribe struct {
	ID shared.IDType

	UserID   shared.IDType
	CourseID shared.IDType

	Level SubscribeLevel

	CreatedAt time.Time
}

type CourseSubscribeRepository interface {
	GetUserSubscribedCourses(ctx context.Context, userID shared.IDType) ([]CourseSubscribe, error)
	Create(ctx context.Context, subscribe *CourseSubscribe) error
	Delete(ctx context.Context, subscribe *CourseSubscribe) error
}

type CourseSubscribeService interface {
	Subscribe(ctx context.Context, req shared.RequestCtx, cmd CourseSubscribeCommand) error
}
