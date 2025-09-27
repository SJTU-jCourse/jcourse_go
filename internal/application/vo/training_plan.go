package vo

import "jcourse_go/internal/infrastructure/entity"

type TrainingPlanVO struct {
	ID         int64   `json:"id"`
	Degree     string  `json:"degree"`
	Major      string  `json:"major"`
	Department string  `json:"department"`
	EntryYear  string  `json:"entry_year"`
	MajorCode  string  `json:"major_code"`
	TotalYear  int64   `json:"total_year"`
	MinCredits float64 `json:"min_credits"`
	MajorClass string  `json:"major_class"`
	CreatedAt  int64   `json:"created_at"`
	UpdatedAt  int64   `json:"updated_at"`
}

type TrainingPlanDetailVO struct {
	ID         int64                      `json:"id"`
	Degree     string                     `json:"degree"`
	Major      string                     `json:"major"`
	Department string                     `json:"department"`
	EntryYear  string                     `json:"entry_year"`
	MajorCode  string                     `json:"major_code"`
	TotalYear  int64                      `json:"total_year"`
	MinCredits float64                    `json:"min_credits"`
	MajorClass string                     `json:"major_class"`
	Curricula  []TrainingPlanCurriculumVO `json:"curricula"`
	CreatedAt  int64                      `json:"created_at"`
	UpdatedAt  int64                      `json:"updated_at"`
}

type CurriculumVO struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Credit float64 `json:"credit"`
}

type TrainingPlanCurriculumVO struct {
	Curriculum      CurriculumVO `json:"curriculum"`
	SuggestSemester string       `json:"suggest_semester"`
	Category        string       `json:"category"`
}

func NewCurriculumVOFromEntity(c *entity.Curriculum) CurriculumVO {
	return CurriculumVO{
		Code:   c.Code,
		Name:   c.Name,
		Credit: c.Credit,
	}
}

func NewTrainingPlanVOFromEntity(t *entity.TrainingPlan) TrainingPlanVO {
	return TrainingPlanVO{
		ID:         t.ID,
		Degree:     t.Degree,
		Major:      t.Major,
		Department: t.Department,
		EntryYear:  t.EntryYear,
		MajorCode:  t.MajorCode,
		TotalYear:  t.TotalYear,
		MinCredits: t.MinCredits,
		MajorClass: t.MajorClass,
		CreatedAt:  t.CreatedAt.Unix(),
		UpdatedAt:  t.UpdatedAt.Unix(),
	}
}

func NewTrainingPlanDetailVOFromEntity(t *entity.TrainingPlan) TrainingPlanDetailVO {
	planVO := TrainingPlanDetailVO{
		ID:         t.ID,
		Degree:     t.Degree,
		Major:      t.Major,
		Department: t.Department,
		EntryYear:  t.EntryYear,
		MajorCode:  t.MajorCode,
		TotalYear:  t.TotalYear,
		MinCredits: t.MinCredits,
		MajorClass: t.MajorClass,
		Curricula:  make([]TrainingPlanCurriculumVO, 0),
		CreatedAt:  t.CreatedAt.Unix(),
		UpdatedAt:  t.UpdatedAt.Unix(),
	}

	for _, c := range t.Curriculums {
		planVO.Curricula = append(planVO.Curricula, NewTrainingPlanCurriculumVOFromEntity(&c))
	}

	return planVO
}

func NewTrainingPlanCurriculumVOFromEntity(t *entity.TrainingPlanCurriculum) TrainingPlanCurriculumVO {
	return TrainingPlanCurriculumVO{
		Curriculum:      NewCurriculumVOFromEntity(t.Curriculum),
		SuggestSemester: t.SuggestSemester,
		Category:        t.Category,
	}
}
