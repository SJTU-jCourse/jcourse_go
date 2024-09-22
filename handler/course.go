package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"jcourse_go/constant"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/service"
)

func GetCourseDetailHandler(c *gin.Context) {
	var request dto.CourseDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "参数错误"})
		return
	}

	course, err := service.GetCourseDetail(c, request.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	c.JSON(http.StatusOK, course)
}

func convertCourseListFilter(request dto.CourseListRequest) model.CourseListFilter {
	filter := model.CourseListFilter{
		Page:          request.Page,
		PageSize:      request.PageSize,
		Categories:    make([]string, 0),
		Departments:   make([]string, 0),
		Credits:       make([]float64, 0),
		Code:          request.Code,
		MainTeacherID: request.MainTeacherID,
	}

	categories := strings.Split(request.Categories, ",")
	for _, category := range categories {
		if category == "" {
			continue
		}
		filter.Categories = append(filter.Categories, category)
	}

	departments := strings.Split(request.Departments, ",")
	for _, department := range departments {
		if department == "" {
			continue
		}
		filter.Departments = append(filter.Departments, department)
	}

	creditsStr := strings.Split(request.Credits, ",")
	for _, creditStr := range creditsStr {
		credit, err := strconv.ParseFloat(creditStr, 64)
		if err != nil {
			continue
		}
		filter.Credits = append(filter.Credits, credit)
	}

	return filter
}

func GetCourseListHandler(c *gin.Context) {
	request := dto.CourseListRequest{
		Page:     constant.DefaultPage,
		PageSize: constant.DefaultPageSize,
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	filter := convertCourseListFilter(request)

	courses, err := service.GetCourseList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	total, err := service.GetCourseCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	resp := dto.CourseListResponse{
		Total:    total,
		Data:     courses,
		Page:     request.Page,
		PageSize: int64(len(courses)),
	}
	c.JSON(http.StatusOK, resp)
}

func GetBaseCourse(c *gin.Context) {
	var request dto.BaseCourseDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	baseCourse, err := service.GetBaseCourse(c, request.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, baseCourse)
}

func GetSuggestedCourseHandler(c *gin.Context) {

}

func WatchCourseHandler(c *gin.Context) {}

func UnWatchCourseHandler(c *gin.Context) {}
