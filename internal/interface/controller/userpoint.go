package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/interface/dto"
)

type UserPointController struct {
	userPointQuery query.UserPointQueryService
}

func NewUserPointController(userPointQuery query.UserPointQueryService) *UserPointController {
	return &UserPointController{
		userPointQuery: userPointQuery,
	}
}

func (c *UserPointController) GetUserPoints(ctx *gin.Context) {
	userID := shared.IDType(0)

	totalValue, userPoints, err := c.userPointQuery.GetUserPoint(ctx, userID)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}

	response := map[string]interface{}{
		"total_value": totalValue,
		"user_points": userPoints,
	}
	dto.WriteDataResponse(ctx, response)
}
