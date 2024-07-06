package po

import "gorm.io/gorm"

type BaseCoursePO struct {
	gorm.Model
	Code   string  `gorm:"index;uniqueIndex"`
	Name   string  `gorm:"index"`
	Credit float64 `gorm:"index"`
}

func (po *BaseCoursePO) TableName() string {
	return "base_courses"
}

type TeacherPO struct {
	gorm.Model
	Name       string `gorm:"index"`
	Code       string `gorm:"index;uniqueIndex"`
	Department string `gorm:"index"`
	Title      string
	Pinyin     string `gorm:"index"`
	PinyinAbbr string `gorm:"index"`
}

func (po *TeacherPO) TableName() string {
	return "teachers"
}

type CoursePO struct {
	gorm.Model
	Code            string  `gorm:"index;index:uniq_course,unique"`
	Name            string  `gorm:"index"`
	Credit          float64 `gorm:"index"`
	MainTeacherID   int64   `gorm:"index;index:uniq_course,unique"`
	MainTeacherName string  `gorm:"index"`
	Department      string  `gorm:"index;index:uniq_course,unique"`
}

func (po *CoursePO) TableName() string {
	return "courses"
}

type CourseCategoryPO struct {
	gorm.Model
	CourseID int64  `gorm:"index;index:uniq_offered_course_category,unique"`
	Category string `gorm:"index;index:uniq_offered_course_category,unique"`
}

func (po *CourseCategoryPO) TableName() string {
	return "course_categories"
}

type OfferedCoursePO struct {
	gorm.Model
	CourseID      int64  `gorm:"index;index:uniq_offered_course,unique"`
	MainTeacherID int64  `gorm:"index"`
	Semester      string `gorm:"index;index:uniq_offered_course,unique"`
	Language      string `gorm:"index"`
	Grade         string `gorm:"index"`
}

func (po *OfferedCoursePO) TableName() string {
	return "offered_courses"
}

type OfferedCourseTeacherPO struct {
	gorm.Model
	CourseID        int64  `gorm:"index"`
	OfferedCourseID int64  `gorm:"index;index:uniq_offered_course_teacher,unique"`
	MainTeacherID   int64  `gorm:"index"`
	TeacherID       int64  `gorm:"index;index:uniq_offered_course_teacher,unique"`
	TeacherName     string `gorm:"index"`
}

func (po *OfferedCourseTeacherPO) TableName() string {
	return "offered_courses_teachers"
}

type TrainingPlanPO struct {
	gorm.Model
	Degree     string `gorm:"index;index:uniq_training_plan,unique"`
	Major      string `gorm:"index;index:uniq_training_plan,unique"`
	Department string `gorm:"index;index:uniq_training_plan,unique"`
	EntryYear  string `gorm:"index;index:uniq_training_plan,unique"`
}

func (po *TrainingPlanPO) TableName() string {
	return "training_plans"
}

type TrainingPlanCoursePO struct {
	gorm.Model
	CourseID       int64 `gorm:"index;index:uniq_training_plan_course,unique"`
	TrainingPlanID int64 `gorm:"index;index:uniq_training_plan_course,unique"`
}

func (po *TrainingPlanCoursePO) TableName() string {
	return "training_plan_courses"
}
