package notification

import (
	"context"
	"time"

	"jcourse_go/internal/domain/shared"
)

type Level int64

const (
	LevelNormal Level = iota
	LevelWatch
	LevelIgnore
)

type CourseNotification struct {
	ID shared.IDType

	UserID   shared.IDType
	CourseID shared.IDType

	Level Level

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CourseNotificationRepository interface {
	Save(ctx context.Context, notification *CourseNotification) error
	Delete(ctx context.Context, notification *CourseNotification) error
}
