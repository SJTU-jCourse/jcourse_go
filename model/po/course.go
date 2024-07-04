package po

import "gorm.io/gorm"

type CoursePO struct {
	gorm.Model
	Code            string  `gorm:"index"`
	Name            string  `gorm:"index"`
	Credit          float64 `gorm:"index"`
	BaseCourseID    int64   `gorm:"index;index:uniq_course,unique"`
	MainTeacherID   int64   `gorm:"index;index:uniq_course,unique"`
	MainTeacherName string  `gorm:"index"`
}

func (po *CoursePO) TableName() string {
	return "courses"
}

type OfferedCoursePO struct {
	gorm.Model
	BaseCourseID  int64  `gorm:"index"`
	CourseID      int64  `gorm:"index;index:uniq_offered_course,unique"`
	MainTeacherID int64  `gorm:"index;index:uniq_offered_course,unique"`
	Semester      string `gorm:"index;index:uniq_offered_course,unique"`
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
	OfferedCourseID int64  `gorm:"index;index:uniq_offered_course_teacher,unique"`
	MainTeacherID   int64  `gorm:"index"`
	TeacherID       int64  `gorm:"index;index:uniq_offered_course_teacher,unique"`
	TeacherName     string `gorm:"index"`
	Semester        string `gorm:"index;index:uniq_offered_course_teacher,unique"`
}

func (po *OfferedCourseTeacherPO) TableName() string {
	return "offered_courses_teachers"
}

type CourseCategoryPO struct {
	gorm.Model
	BaseCourseID  int64  `gorm:"index"`
	CourseID      int64  `gorm:"index;index:uniq_offered_course_category,unique"`
	MainTeacherID int64  `gorm:"index"`
	Category      string `gorm:"index;index:uniq_offered_course_category,unique"`
}

func (po *CourseCategoryPO) TableName() string {
	return "course_categories"
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
