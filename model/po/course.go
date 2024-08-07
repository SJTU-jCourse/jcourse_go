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
	// 联合唯一索引
	Degree     string  `gorm:"index;index:uniq_training_plan,unique"`
	Major      string  `gorm:"index;index:uniq_training_plan,unique"`
	Department string  `gorm:"index;index:uniq_training_plan,unique"`
	EntryYear  string  `gorm:"index;index:uniq_training_plan,unique"` //==Grade,年级
	MajorCode  string  `gorm:"index;index:uniq_training_plan,unique"`
	TotalYear  int     `gorm:"index;index:uniq_training_plan,unique"`
	MinPoints  float64 `gorm:"index;index:uniq_training_plan,unique"`
	MajorClass string  `gorm:"index;index:uniq_training_plan,unique"` // 专业类
}

func (po *TrainingPlanPO) TableName() string {
	return "training_plans"
}

type TrainingPlanCoursePO struct {
	gorm.Model
	CourseID        int64  `gorm:"index;index:uniq_training_plan_course,unique"`
	TrainingPlanID  int64  `gorm:"index;index:uniq_training_plan_course,unique"`
	SuggestYear     int64  `gorm:"index;index:uniq_training_plan_course,unique"`
	SuggestSemester int64  `gorm:"index;index:uniq_training_plan_course,unique"`
	Department      string `gorm:"index;"`
}

func (po *TrainingPlanCoursePO) TableName() string {
	return "training_plan_courses"
}

type TrainingPlanRatePO struct {
	gorm.Model
	UserID         int64 `gorm:"index;index:uniq_training_plan_rate,unique"`
	TrainingPlanID int64 `gorm:"index;index:uniq_training_plan_rate,unique"`
	Rate           int64 `gorm:"index"`
}
type TrainingPlanRateInfoPO struct {
	Average float64
	Count   int64
	Rates    []TrainingPlanRatePO
}

func (po *TrainingPlanRatePO) TableName() string {
	return "training_plan_rates"
}

type CourseReviewInfo struct {
	CourseID int64
	Average  float64
	Count    int64
}
