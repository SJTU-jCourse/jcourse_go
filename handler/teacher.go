package handler

import (
	"net/http"
	"strings"

	"jcourse_go/constant"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"

	"jcourse_go/service"

	"github.com/gin-gonic/gin"
)

func convertTeacherListFilter(request dto.TeacherListRequest) model.TeacherFilterForQuery {
	filter := model.TeacherFilterForQuery{
		PaginationFilterForQuery: request.PaginationFilterForQuery,
		Name:                     request.Name,
		Departments:              make([]string, 0),
		Titles:                   make([]string, 0),
		Pinyin:                   request.Pinyin,
		PinyinAbbr:               request.PinyinAbbr,
	}

	departments := strings.Split(request.Departments, ",")
	for _, d := range departments {
		if d == "" {
			continue
		}
		filter.Departments = append(filter.Departments, d)
	}

	titles := strings.Split(request.Titles, ",")
	for _, t := range titles {
		if t == "" {
			continue
		}
		filter.Titles = append(filter.Titles, t)
	}

	return filter
}

func GetTeacherListHandler(c *gin.Context) {
	request := dto.TeacherListRequest{
		PaginationFilterForQuery: model.PaginationFilterForQuery{
			Page:     constant.DefaultPage,
			PageSize: constant.DefaultPageSize,
		},
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	filter := convertTeacherListFilter(request)

	teachers, err := service.SearchTeacherList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	total, err := service.GetTeacherCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "查无此人。"})
		return
	}

	resp := dto.TeacherListResponse{
		Total:    total,
		Data:     teachers,
		Page:     request.Page,
		PageSize: int64(len(teachers)),
	}
	c.JSON(http.StatusOK, resp)
}

func GetTeacherDetailHandler(c *gin.Context) {
	var request dto.TeacherDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "参数错误"})
	}

	teacher, err := service.GetTeacherDetail(c, request.TeacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, teacher)
}

func SearchTeacherListHandler(c *gin.Context) {
	var request dto.TeacherListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	filter := model.TeacherFilterForQuery{
		Name:                     request.Name,
		Pinyin:                   request.Pinyin,
		PinyinAbbr:               request.PinyinAbbr,
		PaginationFilterForQuery: request.PaginationFilterForQuery,
	}

	teachers, err := service.SearchTeacherList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	total, err := service.GetTeacherCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "查无此人。"})
		return
	}

	resp := dto.TeacherListResponse{
		Total:    total,
		Data:     teachers,
		Page:     request.Page,
		PageSize: int64(len(teachers)),
	}
	c.JSON(http.StatusOK, resp)
}

func GetTeacherFilter(c *gin.Context) {
	filter, err := service.GetTeacherFilter(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, filter)
}

func CreateTeacherHandler(c *gin.Context) {}

func UpdateTeacherHandler(c *gin.Context) {}
