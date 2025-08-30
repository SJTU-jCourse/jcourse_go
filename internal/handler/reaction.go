package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/middleware"
	dto2 "jcourse_go/internal/model/dto"
	"jcourse_go/internal/service"
)

func CreateReviewReactionHandler(c *gin.Context) {
	var request dto2.CreateReviewReactionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto2.BaseResponse{Message: "用户未登录！"})
		return
	}

	id, err := service.CreateReviewReaction(c, request, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误"})
	}
	c.JSON(http.StatusOK, dto2.CreateReviewReactionResponse{ReactionID: id}) // nolint: gosimple
}

func DeleteReviewReactionHandler(c *gin.Context) {
	var request dto2.DeleteReviewReactionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto2.BaseResponse{Message: "用户未登录！"})
		return
	}

	err := service.DeleteReviewReaction(c, user, request.ReactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误"})
	}
	c.JSON(http.StatusOK, dto2.BaseResponse{Message: "删除成功"})
}
