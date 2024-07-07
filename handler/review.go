package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/service"
)

func GetSuggestedReviewHandler(c *gin.Context) {}

func GetReviewDetailHandler(c *gin.Context) {}

func GetReviewListHandler(c *gin.Context) {
	var request dto.ReviewListRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	filter := domain.ReviewFilter{
		Page:     request.Page,
		PageSize: request.PageSize,
	}

	reviews, err := service.GetReviewList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	total, err := service.GetReviewCount(c, filter)

	response := dto.ReviewListResponse{
		Page:     request.Page,
		PageSize: request.PageSize,
		Total:    total,
		Data:     converter.ConvertReviewDomainToListDTO(reviews, true),
	}
	c.JSON(http.StatusOK, response)
}

func CreateReviewHandler(c *gin.Context) {}

func UpdateReviewHandler(c *gin.Context) {}

func DeleteReviewHandler(c *gin.Context) {}

func GetReviewListForCourseHandler(c *gin.Context) {}
