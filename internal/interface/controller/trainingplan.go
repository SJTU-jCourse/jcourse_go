package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/domain/course"
)

type TrainingPlanController struct {
	trainingPlanRepo course.TrainingPlanRepository
	curriculumRepo   course.CurriculumRepository
}

func NewTrainingPlanController(
	trainingPlanRepo course.TrainingPlanRepository,
	curriculumRepo course.CurriculumRepository,
) *TrainingPlanController {
	return &TrainingPlanController{
		trainingPlanRepo: trainingPlanRepo,
		curriculumRepo:   curriculumRepo,
	}
}

func (c *TrainingPlanController) GetTrainingPlanList(ctx *gin.Context) {

}

func (c *TrainingPlanController) GetTrainingPlanDetail(ctx *gin.Context) {

}
