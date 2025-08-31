package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/interface/middleware"
	"jcourse_go/internal/service"
)

func CreateRatingHandler(c *gin.Context) {
	var request dto.RatingDTO
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
		return
	}
	user := middleware.GetCurrentUser(c)
	err := service.CreateRating(c, user.ID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetUserRatingHandler(c *gin.Context) {}
