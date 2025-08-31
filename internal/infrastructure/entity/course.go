package entity

import (
	"time"
)

type BaseCourse struct {
	Code      string    `gorm:"index;uniqueIndex"`
	Name      string    `gorm:"index"`
	Credit    float64   `gorm:"index"`
	CreatedAt time.Time `gorm:"index"`
}

func (s *BaseCourse) TableName() string {
	return "base_courses"
}

type Course struct {
	ID int64 `gorm:"primarykey"`

	Code          string  `gorm:"index;index:uniq_course,unique"`
	Name          string  `gorm:"index"`
	Credit        float64 `gorm:"index"`
	MainTeacherID int64   `gorm:"index;index:uniq_course,unique"`

	SearchIndex SearchIndex `gorm:"->:false;<-"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (po *Course) TableName() string {
	return "courses"
}

type OfferedCoursePO struct {
	ID int64 `gorm:"primarykey"`

	CourseID      int64  `gorm:"index;index:uniq_offered_course,unique"`
	MainTeacherID int64  `gorm:"index"`
	Semester      string `gorm:"index;index:uniq_offered_course,unique"`
	Language      string `gorm:"index"`
	Grade         string `gorm:"index"`
	Department    string `gorm:"index;index:uniq_course,unique"`

	CreatedAt time.Time
}

func (po *OfferedCoursePO) TableName() string {
	return "offered_courses"
}

type OfferedCourseCategoryPO struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time

	CourseID int64   `gorm:"index;index:uniq_offered_course_category,unique"`
	Course   *Course `gorm:"constraint:OnDelete:CASCADE;"`
	Category string  `gorm:"index;index:uniq_offered_course_category,unique"`
}

func (po *OfferedCourseCategoryPO) TableName() string {
	return "offered_course_categories"
}

type OfferedCourseTeacherPO struct {
	ID int64 `gorm:"primarykey"`

	CourseID        int64 `gorm:"index"`
	OfferedCourseID int64 `gorm:"index;index:uniq_offered_course_teacher,unique"`
	MainTeacherID   int64 `gorm:"index"`
	TeacherID       int64 `gorm:"index;index:uniq_offered_course_teacher,unique"`

	CreatedAt time.Time
}

func (po *OfferedCourseTeacherPO) TableName() string {
	return "offered_courses_teachers"
}
