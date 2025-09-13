package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
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

}
