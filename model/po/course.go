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

	RatingCount int64   `gorm:"index;default:0;not null"`
	RatingAvg   float64 `gorm:"index;default:0;not null"`

	SearchIndex SearchIndex `gorm:"->:false;<-"`
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

type CourseReviewInfo struct {
	CourseID int64
	Average  float64
	Count    int64
}
