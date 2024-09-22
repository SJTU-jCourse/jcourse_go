package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/middleware"
	"jcourse_go/model/dto"
	"jcourse_go/service"
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
