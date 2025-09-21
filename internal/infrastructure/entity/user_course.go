package entity

import "time"

type CourseNotification struct {
	ID        int64
	CourseID  int64
	UserID    int64
	Level     int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (po *CourseNotification) TableName() string {
	return "course_notification"
}

type UserCourseEnrollment struct {
	ID        int64
	CourseID  int64
	UserID    int64
	Semester  string
	CreatedAt time.Time
}

func (po *UserCourseEnrollment) TableName() string {
	return "user_course_enrollment"
}
