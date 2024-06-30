package po

import "gorm.io/gorm"

type CoursePO struct {
	gorm.Model
	BaseCourseID  int64 `gorm:"index"`
	MainTeacherID int64 `gorm:"index"`
}

func (po *CoursePO) TableName() string {
	return "courses"
}

type OfferedCoursePO struct {
	gorm.Model
	CoursePO
	Semester   string `gorm:"index"`
	Department string `gorm:"index"`
	Location   string
	Language   string `gorm:"index"`
	Grade      string `gorm:"index"`
}

func (po *OfferedCoursePO) TableName() string {
	return "offered_courses"
}

type OfferedCourseTeacherPO struct {
	gorm.Model
	OfferedCourseID int64 `gorm:"index"`
	TeacherID       int64 `gorm:"index"`
}

func (po *OfferedCourseTeacherPO) TableName() string {
	return "offered_courses_teachers"
}

type CourseCategoryPO struct {
	gorm.Model
	OfferedCourseID int64  `gorm:"index"`
	Category        string `gorm:"index"`
}

func (po *CourseCategoryPO) TableName() string {
	return "offered_course_categories"
}

type TrainingPlanPO struct {
	gorm.Model
	Degree     string
	Major      string
	Department string
	EntryYear  string

	Version int64 `gorm:"index"`
}

func (po *TrainingPlanPO) TableName() string {
	return "training_plans"
}

type TrainingPlanCoursePO struct {
	gorm.Model
	CourseID       int64
	TrainingPlanID int64
}

func (po *TrainingPlanCoursePO) TableName() string {
	return "training_plan_courses"
}

type BaseCoursePO struct {
	gorm.Model
	Code   string  `gorm:"index"`
	Name   string  `gorm:"index"`
	Credit float32 `gorm:"index"`
}

func (po *BaseCoursePO) TableName() string {
	return "base_courses"
}

type TeacherPO struct {
	gorm.Model
	Name       string `gorm:"index"`
	Department string `gorm:"index"`
	Title      string
}

func (po *TeacherPO) TableName() string {
	return "teachers"
}
