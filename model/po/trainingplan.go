package po

import "gorm.io/gorm"

type TrainingPlanPO struct {
	gorm.Model
	Degree     string  `gorm:"index;index:uniq_training_plan,unique"`
	Major      string  `gorm:"index;index:uniq_training_plan,unique"`
	Department string  `gorm:"index;index:uniq_training_plan,unique"`
	EntryYear  string  `gorm:"index;index:uniq_training_plan,unique"` // == grade
	MajorCode  string  `gorm:"index;index:uniq_training_plan,unique"`
	TotalYear  int     `gorm:"index;index:uniq_training_plan,unique"`
	MinCredits float64 `gorm:"index;index:uniq_training_plan,unique"`
	MajorClass string  `gorm:"index;index:uniq_training_plan,unique"` // the class of major

	SearchIndex SearchIndex `gorm:"index:idx_search, type:gin"`
}

func (po *TrainingPlanPO) TableName() string {
	return "training_plans"
}

type TrainingPlanCoursePO struct {
	gorm.Model
	CourseID       int64 `gorm:"index;index:uniq_training_plan_course,unique"`
	TrainingPlanID int64 `gorm:"index;index:uniq_training_plan_course,unique"`
	// SuggestSemester:year+semester e.g. 2023-2024-2
	SuggestSemester string `gorm:"index;index:uniq_training_plan_course,unique"`
	Department      string `gorm:"index;"`
}

func (po *TrainingPlanCoursePO) TableName() string {
	return "training_plan_courses"
}
