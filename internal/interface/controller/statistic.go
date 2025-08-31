package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/domain/statistic"
)

type StatisticController struct {
	statisticRepo statistic.StatisticRepository
}

func NewStatisticController(
	statisticRepo statistic.StatisticRepository,
) *StatisticController {
	return &StatisticController{
		statisticRepo: statisticRepo,
	}
}

func (c *StatisticController) GetLatestStatistics(ctx *gin.Context) {

}
