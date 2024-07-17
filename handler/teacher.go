package handler

import (
	"jcourse_go/constant"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"net/http"

	"jcourse_go/service"

	"github.com/gin-gonic/gin"
)

func convertTeacherListFilter(request dto.TeacherListRequest) domain.TeacherListFilter {
	filter := domain.TeacherListFilter{
		Page:       request.Page,
		PageSize:   request.PageSize,
		Name:       request.Name,
		Department: request.Department,
		Pinyin:     request.Pinyin,
		PinyinAbbr: request.PinyinAbbr,
	}

	return filter
}

func GetTeacherListHandler(c *gin.Context) {
	request := dto.TeacherListRequest{
		Page:     constant.DefaultPage,
		PageSize: constant.DefaultPageSize,
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
		Data:     converter.ConvertTeacherListDomainToDTO(teachers),
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
	c.JSON(http.StatusOK, converter.ConvertTeacherDomainToDTO(*teacher))
}

func SearchTeacherListHandler(c *gin.Context) {
	var request dto.TeacherListRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	filter := domain.TeacherListFilter{
		Name:       request.Name,
		Department: request.Department,
		Pinyin:     request.Pinyin,
		PinyinAbbr: request.PinyinAbbr,
		Page:       request.Page,
		PageSize:   request.PageSize,
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
		Data:     converter.ConvertTeacherListDomainToDTO(teachers),
		Page:     request.Page,
		PageSize: int64(len(teachers)),
	}
	c.JSON(http.StatusOK, resp)
}

func CreateTeacherHandler(c *gin.Context) {}

func UpdateTeacherHandler(c *gin.Context) {}
