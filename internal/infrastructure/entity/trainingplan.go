package entity

import (
	"time"
)

type TrainingPlanPO struct {
	ID int64 `gorm:"primarykey"`

	Degree     string  `gorm:"index;index:uniq_training_plan,unique"` // 学位层次（e.g. 本科）
	Major      string  `gorm:"index;index:uniq_training_plan,unique"`
	Department string  `gorm:"index;index:uniq_training_plan,unique"`
	EntryYear  string  `gorm:"index;index:uniq_training_plan,unique"` // 年级（入学年份）
	MajorCode  string  `gorm:"index;"`                                // 专业代码
	TotalYear  int64   `gorm:"index;"`                                // 学制（年限）
	MinCredits float64 `gorm:"index;"`                                // 最小学分
	MajorClass string  `gorm:"index;"`                                // 学位类型（e.g. 工学）

	SearchIndex SearchIndex `gorm:"->:false;<-"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
}

func (po *TrainingPlanPO) TableName() string {
	return "training_plans"
}

type TrainingPlanCoursePO struct {
	ID int64 `gorm:"primarykey"`

	BaseCourseID    int64  `gorm:"index;index:uniq_training_plan_course,unique"`
	TrainingPlanID  int64  `gorm:"index;index:uniq_training_plan_course,unique"`
	SuggestSemester string `gorm:"index;index:uniq_training_plan_course,unique"` // e.g. 2023-2024-2
	Category        string `gorm:"index;"`

	CreatedAt time.Time `gorm:"index"`
}

func (po *TrainingPlanCoursePO) TableName() string {
	return "training_plan_courses"
}
