package handler

import (
	"fmt"
	"net/http"

	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/service"

	"github.com/gin-gonic/gin"
)

func GetTrainingPlanListHandler(c *gin.Context) {
	var request dto.TrainingPlanListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	filter := domain.TrainingPlanFilter{
		Page:     int64(request.Page),
		PageSize: int64(request.PageSize),
	}
	trainingPlanList, err := service.SearchTrainingPlanList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}
	data := converter.ConvertTrainingPlanDomainListToDTO(trainingPlanList)
	response := dto.TrainingPlanListResponse{
		Page:     int64(request.Page),
		PageSize: int64(request.PageSize),
		Total:    int64(len(data)),
		Data:     data,
	}
	c.JSON(http.StatusOK, response)
}

func GetTrainingPlanHandler(c *gin.Context) {
	var request dto.TrainingPlanDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "参数错误"})
		return
	}
	trainingPlan, err := service.GetTrainingPlanDetail(c, request.TrainingPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	response := converter.ConvertTrainingPlanDomainToDTO(*trainingPlan)
	c.JSON(http.StatusOK, response)
}

func SearchTrainingPlanHandler(c *gin.Context) {
	var request dto.TrainingPlanListQueryRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	filter := domain.TrainingPlanFilter{
		Major:      request.MajorName,
		EntryYear:  fmt.Sprintf("%d", request.EntryYear),
		Department: request.Department,
		Page:       int64(request.Page),
		PageSize:   int64(request.PageSize),
	}
	trainingPlanList, err := service.SearchTrainingPlanList(c, filter)
	count := service.GetTrainingPlanCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}
	data := converter.ConvertTrainingPlanDomainListToDTO(trainingPlanList)
	response := dto.TrainingPlanListResponse{
		Page:     int64(request.Page),
		PageSize: int64(len(data)),
		Total:    count,
		Data:     data,
	}
	c.JSON(http.StatusOK, response)
}
