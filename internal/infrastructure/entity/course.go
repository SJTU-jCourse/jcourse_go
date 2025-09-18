package entity

import (
	"time"
)

type Course struct {
	ID int64 `gorm:"primaryKey"`

	Code   string  `gorm:"index;index:uniq_course,unique"`
	Name   string  `gorm:"index"`
	Credit float64 `gorm:"index"`

	MainTeacherID int64 `gorm:"index;index:uniq_course,unique"`
	MainTeacher   *Teacher

	Offerings []*CourseOffering

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (po *Course) TableName() string {
	return "course"
}

type CourseOffering struct {
	ID int64 `gorm:"primaryKey"`

	CourseID      int64 `gorm:"index;index:uniq_offering,unique"`
	Course        *Course
	MainTeacherID int64 `gorm:"index"`

	Semester   string `gorm:"index;index:uniq_offering,unique"`
	Language   string `gorm:"index"`
	Department string `gorm:"index"`

	Categories   []CourseOfferingCategory
	TeacherGroup []CourseOfferingTeacher

	CreatedAt time.Time
}

func (po *CourseOffering) TableName() string {
	return "course_offering"
}

type CourseOfferingCategory struct {
	ID               int64 `gorm:"primaryKey"`
	CourseOfferingID int64 `gorm:"index:uniq_offering_category,unique"`
	CourseOffering   *CourseOffering
	Category         string `gorm:"index:uniq_offering_category,unique"`
	CourseID         int64  `gorm:"index"`
	Course           *Course
	CreatedAt        time.Time
}

func (po *CourseOfferingCategory) TableName() string {
	return "course_offering_category"
}

type CourseOfferingTeacher struct {
	ID int64 `gorm:"primaryKey"`

	CourseOfferingID int64 `gorm:"index:uniq_offering_teacher,unique"`
	CourseOffering   *CourseOffering
	TeacherID        int64 `gorm:"index:uniq_offering_teacher,unique"`
	Teacher          *Teacher
	CourseID         int64 `gorm:"index"`
	Course           *Course
	CreatedAt        time.Time
}

func (po *CourseOfferingTeacher) TableName() string {
	return "course_offering_teacher"
}

type CourseNotification struct {
	ID        int64
	CourseID  int64
	UserID    int64
	Level     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (po *CourseNotification) TableName() string {
	return "course_notification"
}
