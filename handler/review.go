package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"jcourse_go/middleware"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/service"
)

func GetSuggestedReviewHandler(c *gin.Context) {}

func GetReviewDetailHandler(c *gin.Context) {
	var request dto.ReviewDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "参数错误"})
		return
	}

	reviews, err := service.GetReviewList(c, domain.ReviewFilter{ReviewID: request.ReviewID})
	if err != nil || len(reviews) == 0 {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	reviewDTO := converter.ConvertReviewDomainToDTO(reviews[0], true)
	c.JSON(http.StatusOK, reviewDTO)
}

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

func CreateReviewHandler(c *gin.Context) {
	var request dto.UpdateReviewDTO
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	user := middleware.GetUser(c)
	reviewID, err := service.CreateReview(c, request, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	c.JSON(http.StatusOK, dto.CreateReviewResponse{ReviewID: reviewID})
}

func UpdateReviewHandler(c *gin.Context) {
	var request dto.UpdateReviewRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "参数错误"})
		return
	}

	var reviewDTO dto.UpdateReviewDTO
	if err := c.ShouldBind(&reviewDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	reviewDTO.ID = request.ReviewID
	user := middleware.GetUser(c)
	err := service.UpdateReview(c, reviewDTO, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateReviewResponse{ReviewID: request.ReviewID})
}

func DeleteReviewHandler(c *gin.Context) {
	var request dto.DeleteReviewRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	err := service.DeleteReview(c, request.ReviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, dto.DeleteReviewResponse{ReviewID: request.ReviewID})
}

func GetReviewListForCourseHandler(c *gin.Context) {}
