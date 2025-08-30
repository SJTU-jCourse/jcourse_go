package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/middleware"
	dto2 "jcourse_go/internal/model/dto"
	"jcourse_go/internal/service"
)

func CreateRatingHandler(c *gin.Context) {
	var request dto2.RatingDTO
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	err := service.CreateRating(c, user.ID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetUserRatingHandler(c *gin.Context) {}
