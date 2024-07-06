package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"jcourse_go/constant"
	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/service"
)

func GetCourseDetailHandler(c *gin.Context) {

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
	courses, err := service.GetCourseList(c, request)
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
