package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
)

type TrainingPlanController struct {
	trainingPlanQuery query.TrainingPlanQueryService
}

func NewTrainingPlanController(
	trainingPlanQuery query.TrainingPlanQueryService,
) *TrainingPlanController {
	return &TrainingPlanController{
		trainingPlanQuery: trainingPlanQuery,
	}
}

func (c *TrainingPlanController) GetTrainingPlanList(ctx *gin.Context) {

}

func (c *TrainingPlanController) GetTrainingPlanDetail(ctx *gin.Context) {

}
