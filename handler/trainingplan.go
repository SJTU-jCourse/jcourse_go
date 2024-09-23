package handler

import (
	"net/http"

	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/service"

	"github.com/gin-gonic/gin"
)

func GetTrainingPlanListHandler(c *gin.Context) {
	var request dto.TrainingPlanListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	filter := model.TrainingPlanFilter{
		Page:     int64(request.Page),
		PageSize: int64(request.PageSize),
	}
	trainingPlanList, err := service.SearchTrainingPlanList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}

	response := dto.TrainingPlanListResponse{
		Page:     int64(request.Page),
		PageSize: int64(request.PageSize),
		Total:    int64(len(trainingPlanList)),
		Data:     trainingPlanList,
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
	c.JSON(http.StatusOK, trainingPlan)
}

func SearchTrainingPlanHandler(c *gin.Context) {
	var request dto.TrainingPlanListQueryRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	filter := model.TrainingPlanFilter{
		Major:      request.MajorName,
		EntryYear:  request.EntryYear,
		Department: request.Department,
		Page:       int64(request.Page),
		PageSize:   int64(request.PageSize),
	}
	trainingPlanList, err := service.SearchTrainingPlanList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}
	count, err := service.GetTrainingPlanCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}
	response := dto.TrainingPlanListResponse{
		Page:     int64(request.Page),
		PageSize: int64(len(trainingPlanList)),
		Total:    count,
		Data:     trainingPlanList,
	}
	c.JSON(http.StatusOK, response)
}
