package po

import "gorm.io/gorm"

type CoursePO struct {
	gorm.Model
	Code            string  `gorm:"index"`
	Name            string  `gorm:"index"`
	Credit          float32 `gorm:"index"`
	BaseCourseID    int64   `gorm:"index"`
	MainTeacherID   int64   `gorm:"index"`
	MainTeacherName string  `gorm:"index"`
}

func (po *CoursePO) TableName() string {
	return "courses"
}

type OfferedCoursePO struct {
	gorm.Model
	BaseCourseID  int64  `gorm:"index"`
	CourseID      int64  `gorm:"index"`
	MainTeacherID int64  `gorm:"index"`
	Semester      string `gorm:"index"`
	Department    string `gorm:"index"`
	Language      string `gorm:"index"`
	Grade         string `gorm:"index"`
}

func (po *OfferedCoursePO) TableName() string {
	return "offered_courses"
}

type OfferedCourseTeacherPO struct {
	gorm.Model
	BaseCourseID    int64  `gorm:"index"`
	CourseID        int64  `gorm:"index"`
	OfferedCourseID int64  `gorm:"index"`
	TeacherID       int64  `gorm:"index"`
	TeacherName     string `gorm:"index"`
}

func (po *OfferedCourseTeacherPO) TableName() string {
	return "offered_courses_teachers"
}

type OfferedCourseCategoryPO struct {
	gorm.Model
	BaseCourseID    int64  `gorm:"index"`
	CourseID        int64  `gorm:"index"`
	OfferedCourseID int64  `gorm:"index"`
	MainTeacherID   int64  `gorm:"index"`
	Category        string `gorm:"index"`
}

func (po *OfferedCourseCategoryPO) TableName() string {
	return "offered_course_categories"
}

type TrainingPlanPO struct {
	gorm.Model
	Degree     string `gorm:"index"`
	Major      string `gorm:"index"`
	Department string `gorm:"index"`
	EntryYear  string `gorm:"index"`
}

func (po *TrainingPlanPO) TableName() string {
	return "training_plans"
}

type TrainingPlanCoursePO struct {
	gorm.Model
	CourseID       int64 `gorm:"index"`
	TrainingPlanID int64 `gorm:"index"`
}

func (po *TrainingPlanCoursePO) TableName() string {
	return "training_plan_courses"
}

type BaseCoursePO struct {
	gorm.Model
	Code   string  `gorm:"index"`
	Name   string  `gorm:"index"`
	Credit float64 `gorm:"index"`
}

func (po *BaseCoursePO) TableName() string {
	return "base_courses"
}

type TeacherPO struct {
	gorm.Model
	Name       string `gorm:"index"`
	Code       string `gorm:"index"`
	Department string `gorm:"index"`
	Title      string
	Pinyin     string `gorm:"index"`
	PinyinAbbr string `gorm:"index"`
}

func (po *TeacherPO) TableName() string {
	return "teachers"
}
