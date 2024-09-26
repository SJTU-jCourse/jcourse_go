package handler

import (
	"net/http"
	"strings"

	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/service"

	"github.com/gin-gonic/gin"
)

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
	filter := convTrainingPlanFilter(request)
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

func GetTrainingPlanFilter(c *gin.Context) {
	filter, err := service.GetTrainingPlanFilter(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, filter)
}

func convTrainingPlanFilter(request dto.TrainingPlanListQueryRequest) model.TrainingPlanFilterForQuery {
	filter := model.TrainingPlanFilterForQuery{
		Major:       request.MajorName,
		EntryYears:  make([]string, 0),
		Departments: make([]string, 0),
		Degrees:     make([]string, 0),
		Page:        int64(request.Page),
		PageSize:    int64(request.PageSize),
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
