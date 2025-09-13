package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/interface/middleware"
	"jcourse_go/internal/service"
)

func CreateReviewReactionHandler(c *gin.Context) {
	var request olddto.CreateReviewReactionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, olddto.BaseResponse{Message: "用户未登录！"})
		return
	}

	id, err := service.CreateReviewReaction(c, request, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误"})
	}
	c.JSON(http.StatusOK, olddto.CreateReviewReactionResponse{ReactionID: id}) // nolint: gosimple
}

func DeleteReviewReactionHandler(c *gin.Context) {
	var request olddto.DeleteReviewReactionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, olddto.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, olddto.BaseResponse{Message: "用户未登录！"})
		return
	}

	err := service.DeleteReviewReaction(c, user, request.ReactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, olddto.BaseResponse{Message: "内部错误"})
	}
	c.JSON(http.StatusOK, olddto.BaseResponse{Message: "删除成功"})
}
