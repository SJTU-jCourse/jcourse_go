package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/middleware"
	"jcourse_go/model/dto"
)

func GetCommonInfo(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto.BaseResponse{Message: "用户未登录！"})
		return
	}
	resp := dto.CommonInfoResponse{
		User: *user,
	}
	c.JSON(http.StatusOK, resp)
}
