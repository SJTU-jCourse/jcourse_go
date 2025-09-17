package command

import (
	"context"
	"time"

	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/notification"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/pkg/apperror"
)

type CourseNotificationService interface {
	Change(ctx context.Context, req shared.RequestCtx, cmd notification.CourseNotificationCommand) error
}

type courseNotificationService struct {
	courseRepo       course.CourseRepository
	notificationRepo notification.CourseNotificationRepository
}

func (s *courseNotificationService) Change(ctx context.Context, req shared.RequestCtx, cmd notification.CourseNotificationCommand) error {
	c, err := s.courseRepo.Get(ctx, cmd.CourseID)
	if err != nil {
		return err
	}
	if c == nil {
		return apperror.ErrNotFound
	}
	n := notification.CourseNotification{
		ID:        0,
		UserID:    req.User.UserID,
		CourseID:  cmd.CourseID,
		Level:     cmd.Level,
		UpdatedAt: time.Now(),
	}
	return s.notificationRepo.Save(ctx, &n)
}

func NewCourseNotificationService(
	courseRepo course.CourseRepository,
	notificationRepo notification.CourseNotificationRepository,
) CourseNotificationService {
	return &courseNotificationService{
		courseRepo:       courseRepo,
		notificationRepo: notificationRepo,
	}
}
