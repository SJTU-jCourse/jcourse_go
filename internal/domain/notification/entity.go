package notification

import (
	"context"
	"time"

	"jcourse_go/internal/domain/shared"
)

type NotificationLevel string

const (
	NotificationLevelNormal NotificationLevel = "normal"
	NotificationLevelWatch  NotificationLevel = "watch"
	NotificationLevelIgnore NotificationLevel = "ignore"
)

type CourseNotification struct {
	ID shared.IDType

	UserID   shared.IDType
	CourseID shared.IDType

	Level NotificationLevel

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CourseNotificationRepository interface {
	Save(ctx context.Context, subscribe *CourseNotification) error
	Delete(ctx context.Context, subscribe *CourseNotification) error
}
