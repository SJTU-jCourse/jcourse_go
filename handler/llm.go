package handler

import (
	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VectorizeCourseReviewsHandler(c *gin.Context) {
	var request dto.VectorizeCourseReviewsRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "参数错误"})
		return
	}

	vector, err := service.VectorizeCourseReviews(c, request.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	courseDetailDTO := converter.ConvertVectorToString(vector)

	c.JSON(http.StatusOK, courseDetailDTO)
}

func GetMatchCourses(c *gin.Context) {
	var request dto.GetMatchCourseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	// TODO

}
