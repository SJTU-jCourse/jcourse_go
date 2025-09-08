package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application"
)

type StatisticController struct {
	statisticQuery application.StatisticQueryService
}

func NewStatisticController(
	statisticQuery application.StatisticQueryService,
) *StatisticController {
	return &StatisticController{
		statisticQuery: statisticQuery,
	}
}

func (c *StatisticController) GetLatestStatistics(ctx *gin.Context) {

}
