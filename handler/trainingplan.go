package handler

import (
	"fmt"
	"jcourse_go/middleware"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTrainingPlanListHandler(c *gin.Context) {
	var request dto.TrainingPlanListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	Filter := domain.TrainingPlanFilter{
		Page:     int64(request.Page),
		PageSize: int64(request.PageSize),
	}
	TrainingPlanList, err := service.SearchTrainingPlanList(c, Filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}
	data := converter.ConvertTrainingPlanDomainListToDTO(TrainingPlanList)
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
	TrainingPlan, err := service.GetTrainingPlanDetail(c, request.TrainingPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	response := converter.ConvertTrainingPlanDomainToDTO(*TrainingPlan)
	c.JSON(http.StatusOK, response)
}

func SearchTrainingPlanHandler(c *gin.Context) {
	var request dto.TrainingPlanListQueryRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	Filter := domain.TrainingPlanFilter{
		Major:      request.MajorName,
		EntryYear:  fmt.Sprintf("%d", request.EntryYear),
		Department: request.Department,
		Page:       int64(request.Page),
		PageSize:   int64(request.PageSize),
	}
	TrainingPlanList, err := service.SearchTrainingPlanList(c, Filter)
	count := service.GetTrainingPlanCount(c, Filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
	}
	data := converter.ConvertTrainingPlanDomainListToDTO(TrainingPlanList)
	response := dto.TrainingPlanListResponse{
		Page:     int64(request.Page),
		PageSize: int64(len(data)),
		Total:    count,
		Data:     data,
	}
	c.JSON(http.StatusOK, response)
}

// ATTENTION: without test now
func RateTrainingPlanHandler(c *gin.Context) {
	var request dto.RateTrainingPlanRequest
	userId := middleware.GetUser(c).ID
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "参数错误"})
		return
	}
	err := service.RateTrainingPlan(c, userId, request.TrainingPlanID, request.Rate)
	// HINT: 底层upsert，不会出现重复插入错误
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		println(err)
		return
	}
	response := dto.BaseResponse{Message: "评分成功"}
	c.JSON(http.StatusOK, response)
}