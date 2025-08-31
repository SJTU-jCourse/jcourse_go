package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/interface/middleware"
	"jcourse_go/internal/service"
)

func GetCommonInfo(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{Message: "用户未登录！"})
		return
	}

	settings, err := service.GetClientSettings(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误"})
		return
	}

	resp := dto.CommonInfoResponse{
		User:     *user,
		Settings: settings,
	}
	c.JSON(http.StatusOK, resp)
}
