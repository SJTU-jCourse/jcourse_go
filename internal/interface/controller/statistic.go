package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
)

type StatisticController struct {
	statisticQuery query.StatisticQueryService
}

func NewStatisticController(
	statisticQuery query.StatisticQueryService,
) *StatisticController {
	return &StatisticController{
		statisticQuery: statisticQuery,
	}
}

func (c *StatisticController) GetLatestStatistics(ctx *gin.Context) {

}
