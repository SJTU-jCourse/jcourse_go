package repository

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/notification"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type CourseNotificationRepository struct {
	db *gorm.DB
}

func (r *CourseNotificationRepository) Delete(ctx context.Context, notification *notification.CourseNotification) error {
	if err := r.db.Model(&entity.CourseNotification{}).Where("id = ?", notification.ID).Delete(nil).Error; err != nil {
		return err
	}
	return nil
}

func (r *CourseNotificationRepository) Save(ctx context.Context, notification *notification.CourseNotification) error {
	e := newNotificationEntityFromDomain(notification)
	if err := r.db.Model(&entity.CourseNotification{}).
		Create(&e).Error; err != nil {
		return err
	}
	notification.ID = shared.IDType(e.ID)
	return nil
}

func NewCourseNotificationRepository(db *gorm.DB) notification.CourseNotificationRepository {
	return &CourseNotificationRepository{db: db}
}

func newNotificationEntityFromDomain(n *notification.CourseNotification) entity.CourseNotification {
	return entity.CourseNotification{
		ID:        int64(n.ID),
		CourseID:  int64(n.CourseID),
		UserID:    int64(n.UserID),
		Level:     int64(n.Level),
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
