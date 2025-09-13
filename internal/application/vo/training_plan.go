package vo

import "jcourse_go/internal/infrastructure/entity"

type TrainingPlanVO struct{}

type CurriculumVO struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Credit float64 `json:"credit"`
}

func NewCurriculumVOFromEntity(c *entity.Curriculum) CurriculumVO {
	return CurriculumVO{
		Code:   c.Code,
		Name:   c.Name,
		Credit: c.Credit,
	}
}
