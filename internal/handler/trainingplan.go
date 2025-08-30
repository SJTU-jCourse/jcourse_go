package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	dto2 "jcourse_go/internal/model/dto"
	"jcourse_go/internal/model/model"
	"jcourse_go/internal/service"
)

func GetTrainingPlanHandler(c *gin.Context) {
	var request dto2.TrainingPlanDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	trainingPlan, err := service.GetTrainingPlanDetail(c, request.TrainingPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, trainingPlan)
}

func SearchTrainingPlanHandler(c *gin.Context) {
	var request dto2.TrainingPlanListQueryRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	filter := convTrainingPlanFilter(request)
	trainingPlanList, err := service.SearchTrainingPlanList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
	}
	count, err := service.GetTrainingPlanCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
	}
	response := dto2.TrainingPlanListResponse{
		Page:     int64(request.Page),
		PageSize: int64(len(trainingPlanList)),
		Total:    count,
		Data:     trainingPlanList,
	}
	c.JSON(http.StatusOK, response)
}

func GetTrainingPlanFilter(c *gin.Context) {
	filter, err := service.GetTrainingPlanFilter(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, filter)
}

func convTrainingPlanFilter(request dto2.TrainingPlanListQueryRequest) model.TrainingPlanFilterForQuery {
	filter := model.TrainingPlanFilterForQuery{
		Major:                    request.MajorName,
		EntryYears:               make([]string, 0),
		Departments:              make([]string, 0),
		Degrees:                  make([]string, 0),
		PaginationFilterForQuery: request.PaginationFilterForQuery,
	}

	degrees := strings.Split(request.Degrees, ",")
	for _, degree := range degrees {
		if degree == "" {
			continue
		}
		filter.Degrees = append(filter.Degrees, degree)
	}

	departments := strings.Split(request.Departments, ",")
	for _, department := range departments {
		if department == "" {
			continue
		}
		filter.Departments = append(filter.Departments, department)
	}

	entryYears := strings.Split(request.EntryYears, ",")
	for _, year := range entryYears {
		if year == "" {
			continue
		}
		filter.EntryYears = append(filter.EntryYears, year)
	}
	return filter
}
