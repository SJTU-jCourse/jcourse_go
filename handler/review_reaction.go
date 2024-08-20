package handler

import (
	"jcourse_go/model/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateReviewReactionHandler(c *gin.Context) {
	var request dto.CreateReviewReactionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}

}

func DeleteReviewReactionHandler(c *gin.Context) {}
