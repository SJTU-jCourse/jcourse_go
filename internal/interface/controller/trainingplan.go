package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/interface/dto"
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
	var req course.TrainingPlanListQuery
	if err := ctx.ShouldBind(&req); err != nil {
		dto.WriteBadArgumentResponse(ctx)
		return
	}
	trainingPlans, err := c.trainingPlanQuery.GetTrainingPlanList(ctx, req)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, trainingPlans)
}

func (c *TrainingPlanController) GetTrainingPlanDetail(ctx *gin.Context) {
	trainingPlanIDStr := ctx.Param("trainingPlanID")
	if trainingPlanIDStr == "" {
		dto.WriteBadArgumentResponse(ctx)
		return
	}
	trainingPlanID, err := strconv.Atoi(trainingPlanIDStr)
	if err != nil || trainingPlanID == 0 {
		dto.WriteBadArgumentResponse(ctx)
		return
	}

	trainingPlan, err := c.trainingPlanQuery.GetTrainingPlanDetail(ctx, shared.IDType(trainingPlanID))
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, trainingPlan)
}

func (c *TrainingPlanController) GetTrainingPlanFilter(ctx *gin.Context) {
	filter, err := c.trainingPlanQuery.GetTrainingPlanFilter(ctx)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, filter)
}
