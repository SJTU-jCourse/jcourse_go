package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application"
)

type TrainingPlanController struct {
	trainingPlanQuery application.TrainingPlanQueryService
}

func NewTrainingPlanController(
	trainingPlanQuery application.TrainingPlanQueryService,
) *TrainingPlanController {
	return &TrainingPlanController{
		trainingPlanQuery: trainingPlanQuery,
	}
}

func (c *TrainingPlanController) GetTrainingPlanList(ctx *gin.Context) {

}

func (c *TrainingPlanController) GetTrainingPlanDetail(ctx *gin.Context) {

}
