package handler

import (
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/service"
	"jcourse_go/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetStatisticHandler(c *gin.Context) {
	request := dto.StatisticRequest{
		StartTime:  0,
		EndTime:    0,
		PeriodKeys: []model.PeriodInfoKey{},
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	if request.StartTime > request.EndTime {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	filter := model.StatisticFilter{
		StartDate:      util.FormatDate(time.Unix(request.StartTime, 0)),
		EndDate:        util.FormatDate(time.Unix(request.EndTime, 0)),
		PeriodInfoKeys: request.PeriodKeys,
	}
	dailyInfos, periodInfos, err := service.GetStatistics(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误"})
		return
	}
	c.JSON(http.StatusOK, dto.StatisticResponse{
		DailyInfos:  dailyInfos,
		PeriodInfos: periodInfos,
	})
}
