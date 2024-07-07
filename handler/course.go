package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"jcourse_go/constant"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
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

	courseDetailDTO := converter.ConvertCourseDomainToDetailDTO(*course)

	c.JSON(http.StatusOK, courseDetailDTO)
}

func convertCourseListFilter(request dto.CourseListRequest) domain.CourseListFilter {
	filter := domain.CourseListFilter{
		Page:        request.Page,
		PageSize:    request.PageSize,
		Categories:  make([]string, 0),
		Departments: make([]string, 0),
		Credits:     make([]float64, 0),
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

	courses, err := service.GetCourseList(c, convertCourseListFilter(request))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	total, err := service.GetCourseCount(c)

	resp := dto.CourseListResponse{
		Total:    total,
		Data:     converter.ConvertCourseListDomainToDTO(courses),
		Page:     request.Page,
		PageSize: int64(len(courses)),
	}
	c.JSON(http.StatusOK, resp)
}

func GetSuggestedCourseHandler(c *gin.Context) {

}

func WatchCourseHandler(c *gin.Context) {}

func UnWatchCourseHandler(c *gin.Context) {}
