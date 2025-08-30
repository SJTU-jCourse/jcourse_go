package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/middleware"
	dto2 "jcourse_go/internal/model/dto"
	"jcourse_go/internal/service"
)

func GetCommonInfo(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto2.BaseResponse{Message: "用户未登录！"})
		return
	}

	settings, err := service.GetClientSettings(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "内部错误"})
		return
	}

	resp := dto2.CommonInfoResponse{
		User:     *user,
		Settings: settings,
	}
	c.JSON(http.StatusOK, resp)
}
