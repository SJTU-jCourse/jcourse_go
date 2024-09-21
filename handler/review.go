package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/constant"
	"jcourse_go/middleware"
	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/service"
)

func GetSuggestedReviewHandler(c *gin.Context) {}

func GetReviewDetailHandler(c *gin.Context) {
	var request dto.ReviewDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "参数错误"})
		return
	}

	reviews, err := service.GetReviewList(c, model.ReviewFilter{ReviewID: request.ReviewID})
	if err != nil || len(reviews) == 0 {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	c.JSON(http.StatusOK, reviews[0])
}

func GetReviewListHandler(c *gin.Context) {
	var request = dto.ReviewListRequest{
		Page:     constant.DefaultPage,
		PageSize: constant.DefaultPageSize,
	}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	filter := model.ReviewFilter{
		Page:     request.Page,
		PageSize: request.PageSize,
	}

	reviews, err := service.GetReviewList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	total, err := service.GetReviewCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	converter.RemoveReviewsUserInfo(reviews, true)

	response := dto.ReviewListResponse{
		Page:     request.Page,
		PageSize: request.PageSize,
		Total:    total,
		Data:     reviews,
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

	c.JSON(http.StatusOK, dto.UpdateReviewResponse{ReviewID: request.ReviewID}) // nolint: gosimple
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

	c.JSON(http.StatusOK, dto.DeleteReviewResponse{ReviewID: request.ReviewID}) // nolint: gosimple
}

func GetReviewListForCourseHandler(c *gin.Context) {}

func GetUserReviewsHandler(c *gin.Context) {
	userID, err := getUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "非法用户ID"})
		return
	}

	var request dto.ReviewListRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

	filter := model.ReviewFilter{
		Page:     request.Page,
		PageSize: request.PageSize,
		UserID:   userID,
	}

	reviews, err := service.GetReviewList(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	total, err := service.GetReviewCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}

	response := dto.ReviewListResponse{
		Page:     request.Page,
		PageSize: request.PageSize,
		Total:    total,
		Data:     reviews,
	}
	c.JSON(http.StatusOK, response)
}
