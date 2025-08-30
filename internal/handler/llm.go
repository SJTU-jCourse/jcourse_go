package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dto2 "jcourse_go/internal/model/dto"
	"jcourse_go/internal/service"
)

func OptCourseReviewHandler(c *gin.Context) {
	var request dto2.OptCourseReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}

	response, err := service.OptCourseReview(c, request.CourseName, request.ReviewContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
		return
	}

	c.JSON(http.StatusOK, response)
}
func GetCourseSummaryHandler(c *gin.Context) {
	var request dto2.CourseDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto2.BaseResponse{Message: "参数错误"})
		return
	}

	response, err := service.GetCourseSummary(c, request.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
		return
	}

	c.JSON(http.StatusOK, response)
}
func VectorizeCourseHandler(c *gin.Context) {
	var request dto2.CourseDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto2.BaseResponse{Message: "参数错误"})
		return
	}

	err := service.VectorizeCourse(c, request.CourseID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
		return
	}

	c.JSON(http.StatusOK, dto2.BaseResponse{Message: "向量化成功"})
}

func GetMatchCoursesHandler(c *gin.Context) {
	var request dto2.GetMatchCourseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}

	courses, err := service.GetMatchCourses(c, request.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
		return
	}

	// TODO: 修改Total, Page
	resp := dto2.CourseListResponse{
		Total:    0,
		Data:     courses,
		Page:     0,
		PageSize: int64(len(courses)),
	}
	c.JSON(http.StatusOK, resp)
}
