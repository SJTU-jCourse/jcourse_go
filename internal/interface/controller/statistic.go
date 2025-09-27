package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
	"jcourse_go/internal/interface/dto"
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

func (c *StatisticController) GetStatistics(ctx *gin.Context) {
	date := ctx.Query("date")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if date != "" {
		statistic, err := c.statisticQuery.GetDailyStatistic(ctx, date)
		if err != nil {
			dto.WriteErrorResponse(ctx, err)
			return
		}
		dto.WriteDataResponse(ctx, statistic)
		return
	}

	if startDate != "" && endDate != "" {
		statistics, err := c.statisticQuery.GetRangeStatistic(ctx, startDate, endDate)
		if err != nil {
			dto.WriteErrorResponse(ctx, err)
			return
		}
		dto.WriteDataResponse(ctx, statistics)
		return
	}

	dto.WriteBadArgumentResponse(ctx)
}
