package po

import (
	"time"
)

type BaseCoursePO struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time

	Code   string  `gorm:"index;uniqueIndex"`
	Name   string  `gorm:"index"`
	Credit float64 `gorm:"index"`
}

func (po *BaseCoursePO) TableName() string {
	return "base_courses"
}

type CoursePO struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Code   string  `gorm:"index;index:uniq_course,unique"`
	Name   string  `gorm:"index"`
	Credit float64 `gorm:"index"`

	MainTeacherID   int64      `gorm:"index;index:uniq_course,unique"`
	MainTeacher     *TeacherPO `gorm:"constraint:OnDelete:CASCADE;"`
	MainTeacherName string     `gorm:"index"`
	Department      string     `gorm:"index;index:uniq_course,unique"`

	RatingCount int64   `gorm:"index;default:0;not null"`
	RatingAvg   float64 `gorm:"index;default:0;not null"`

	Categories     []CourseCategoryPO `gorm:"foreignKey:CourseID"`
	OfferedCourses []OfferedCoursePO  `gorm:"foreignKey:CourseID"`

	SearchIndex SearchIndex `gorm:"->:false;<-"`
}

func (po *CoursePO) TableName() string {
	return "courses"
}

type CourseCategoryPO struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time

	CourseID int64     `gorm:"index;index:uniq_offered_course_category,unique"`
	Course   *CoursePO `gorm:"constraint:OnDelete:CASCADE;"`
	Category string    `gorm:"index;index:uniq_offered_course_category,unique"`
}

func (po *CourseCategoryPO) TableName() string {
	return "course_categories"
}

type OfferedCoursePO struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time

	CourseID      int64      `gorm:"index;index:uniq_offered_course,unique"`
	Course        *CoursePO  `gorm:"constraint:OnDelete:CASCADE;"`
	MainTeacherID int64      `gorm:"index"`
	MainTeacher   *TeacherPO `gorm:"constraint:OnDelete:CASCADE;"`
	Semester      string     `gorm:"index;index:uniq_offered_course,unique"`
	Language      string     `gorm:"index"`
	Grade         string     `gorm:"index"`

	OfferedCourseTeacher []OfferedCourseTeacherPO `gorm:"foreignKey:OfferedCourseID"`
}

func (po *OfferedCoursePO) TableName() string {
	return "offered_courses"
}

type OfferedCourseTeacherPO struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time

	CourseID        int64            `gorm:"index"`
	Course          *CoursePO        `gorm:"constraint:OnDelete:CASCADE;"`
	OfferedCourseID int64            `gorm:"index;index:uniq_offered_course_teacher,unique"`
	OfferedCourse   *OfferedCoursePO `gorm:"constraint:OnDelete:CASCADE;"`
	MainTeacherID   int64            `gorm:"index"`
	MainTeacher     *TeacherPO       `gorm:"constraint:OnDelete:CASCADE;"`
	TeacherID       int64            `gorm:"index;index:uniq_offered_course_teacher,unique"`
	Teacher         *TeacherPO       `gorm:"constraint:OnDelete:CASCADE;"`
	TeacherName     string           `gorm:"index"`
}

func (po *OfferedCourseTeacherPO) TableName() string {
	return "offered_courses_teachers"
}

type CourseReviewInfo struct {
	CourseID int64
	Average  float64
	Count    int64
}
