package handler

import (
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllStatisticHandler(c *gin.Context) {
	statistics, err := service.GetAllStatistics(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误"})
		return
	}
	c.JSON(http.StatusOK, dto.StatisticResponse{DailyInfos: statistics})
}

func GetDailyStatisticHandler(c *gin.Context) {
	request := dto.StatisticRequest{
		InDetail: false,
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	if request.DateTime <= 0 {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	if !request.InDetail {
		minimal, err := service.GetDailyStatisticMinimal(c, time.Unix(request.DateTime, 0))
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误"})
			return
		}
		detail := model.DailyInfoDetail{
			DailyInfoMinimal: minimal,
			WAU:              0,
			MAU:              0,
		}
		c.JSON(http.StatusOK, dto.StatisticResponse{DailyInfos: []model.DailyInfoDetail{
			detail,
		}})
		return
	}

	detail, err := service.GetDailyStatisticDetail(c, time.Unix(request.DateTime, 0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误"})
		return
	}
	c.JSON(http.StatusOK, dto.StatisticResponse{DailyInfos: []model.DailyInfoDetail{
		detail,
	}})
}
