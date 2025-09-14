package entity

import (
	"time"
)

type Curriculum struct {
	Code      string    `gorm:"primaryKey"`
	Name      string    `gorm:"index"`
	Credit    float64   `gorm:"index"`
	CreatedAt time.Time `gorm:"index"`
}

func (s *Curriculum) TableName() string {
	return "curriculum"
}

type TrainingPlan struct {
	ID int64 `gorm:"primaryKey"`

	Degree     string  `gorm:"index;index:uniq_training_plan,unique"` // 学位层次（e.g. 本科）
	Major      string  `gorm:"index;index:uniq_training_plan,unique"`
	Department string  `gorm:"index;index:uniq_training_plan,unique"`
	EntryYear  string  `gorm:"index;index:uniq_training_plan,unique"` // 年级（入学年份）
	MajorCode  string  `gorm:"index;"`                                // 专业代码
	TotalYear  int64   `gorm:"index;"`                                // 学制（年限）
	MinCredits float64 `gorm:"index;"`                                // 最小学分
	MajorClass string  `gorm:"index;"`                                // 学位类型（e.g. 工学）

	Curriculums []TrainingPlanCurriculum

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
}

func (po *TrainingPlan) TableName() string {
	return "training_plan"
}

type TrainingPlanCurriculum struct {
	ID int64 `gorm:"primaryKey"`

	CurriculumCode  string `gorm:"index:uniq_training_plan_course,unique"`
	Curriculum      *Curriculum
	TrainingPlanID  int64 `gorm:"index:uniq_training_plan_course,unique"`
	TrainingPlan    *TrainingPlan
	SuggestSemester string `gorm:"index;"` // e.g. 2023-2024-2
	Category        string `gorm:"index;"`

	CreatedAt time.Time `gorm:"index"`
}

func (po *TrainingPlanCurriculum) TableName() string {
	return "training_plan_curriculum"
}
