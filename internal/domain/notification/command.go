package notification

import "jcourse_go/internal/domain/shared"

type CourseNotificationCommand struct {
	CourseID shared.IDType     `json:"course_id"`
	Level    NotificationLevel `json:"level"`
}

func (cmd *CourseNotificationCommand) Validate() error {
	return nil
}
