package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/interface/middleware"
	"jcourse_go/internal/service"
)

func CreateReviewReactionHandler(c *gin.Context) {
	var request dto.CreateReviewReactionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{Message: "用户未登录！"})
		return
	}

	id, err := service.CreateReviewReaction(c, request, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误"})
	}
	c.JSON(http.StatusOK, dto.CreateReviewReactionResponse{ReactionID: id}) // nolint: gosimple
}

func DeleteReviewReactionHandler(c *gin.Context) {
	var request dto.DeleteReviewReactionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{Message: "用户未登录！"})
		return
	}

	err := service.DeleteReviewReaction(c, user, request.ReactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误"})
	}
	c.JSON(http.StatusOK, dto.BaseResponse{Message: "删除成功"})
}
