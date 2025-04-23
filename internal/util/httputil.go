package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/model/dto"
)

func BaseResponse(c *gin.Context, code int, message string) {
	c.JSON(code, dto.BaseResponse{Message: message})
}

func SuccessSimpleResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, dto.BaseResponse{Message: message})
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func WrongParamResponse(c *gin.Context) {
	c.JSON(http.StatusBadRequest, dto.BaseResponse{Message: "参数错误"})
}

func ErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: message})
}
