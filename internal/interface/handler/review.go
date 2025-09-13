package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/constant"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/review"
	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/interface/middleware"
	"jcourse_go/internal/model/converter"
	"jcourse_go/internal/service"
)

func GetSuggestedReviewHandler(c *gin.Context) {}

func GetReviewDetailHandler(c *gin.Context) {
	var request olddto.ReviewDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, olddto.BaseResponse{Message: "参数错误"})
		return
	}

	currentUserID := int64(0)
	user := middleware.GetCurrentUser(c)
	if user != nil {
		currentUserID = user.ID
	}

	reviews, err := service.GetReviewList(c, user, reaction.ReviewFilterForQuery{ReviewID: request.ReviewID})
	if err != nil || len(reviews) == 0 {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误。"})
		return
	}

	converter.RemoveReviewsUserInfo(reviews, currentUserID, true)

	c.JSON(http.StatusOK, reviews[0])
}

func GetReviewListHandler(c *gin.Context) {
	var request = olddto.ReviewListRequest{
		PaginationFilterForQuery: course.PaginationFilterForQuery{
			Page:      constant.DefaultPage,
			PageSize:  constant.DefaultPageSize,
			Order:     "created_at",
			Ascending: false,
		},
	}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "参数错误"})
		return
	}

	filter := reaction.ReviewFilterForQuery{
		PaginationFilterForQuery: request.PaginationFilterForQuery,
		UserID:                   request.UserID,
		CourseID:                 request.CourseID,
		Rating:                   request.Rating,
	}

	// 非本人不可看匿名点评
	currentUserID := int64(0)
	user := middleware.GetCurrentUser(c)
	if request.UserID != 0 && (user == nil || user.ID != request.UserID) {
		filter.ExcludeAnonymous = true
	}
	if user != nil {
		currentUserID = user.ID
	}

	reviews, err := service.GetReviewList(c, user, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误。"})
		return
	}
	total, err := service.GetReviewCount(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误。"})
		return
	}

	converter.RemoveReviewsUserInfo(reviews, currentUserID, true)

	response := olddto.ReviewListResponse{
		Page:     request.Page,
		PageSize: request.PageSize,
		Total:    total,
		Data:     reviews,
	}
	c.JSON(http.StatusOK, response)
}

func CreateReviewHandler(c *gin.Context) {
	var request olddto.UpdateReviewDTO
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "参数错误"})
		return
	}

	user := middleware.GetCurrentUser(c)
	reviewID, err := service.CreateReview(c, request, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误。"})
		return
	}

	c.JSON(http.StatusOK, olddto.CreateReviewResponse{ReviewID: reviewID})
}

func UpdateReviewHandler(c *gin.Context) {
	var request olddto.UpdateReviewRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, olddto.BaseResponse{Message: "参数错误"})
		return
	}

	var reviewDTO olddto.UpdateReviewDTO
	if err := c.ShouldBind(&reviewDTO); err != nil {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "参数错误"})
		return
	}

	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, olddto.BaseResponse{Message: "用户未登录！"})
		return
	}
	reviewDTO.ID = request.ReviewID

	err := service.UpdateReview(c, reviewDTO, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误。"})
		return
	}

	c.JSON(http.StatusOK, olddto.UpdateReviewResponse{ReviewID: request.ReviewID}) // nolint: gosimple
}

func DeleteReviewHandler(c *gin.Context) {
	var request olddto.DeleteReviewRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, olddto.BaseResponse{Message: "用户未登录！"})
		return
	}

	err := service.DeleteReview(c, request.ReviewID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误。"})
		return
	}

	c.JSON(http.StatusOK, olddto.DeleteReviewResponse{ReviewID: request.ReviewID}) // nolint: gosimple
}
