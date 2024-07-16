package handler

import (
	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"net/http"

	"jcourse_go/service"

	"github.com/gin-gonic/gin"
)

func GetTeacherListHandler(c *gin.Context) {}

func GetTeacherDetailHandler(c *gin.Context) {
	var request dto.TeacherDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusNotFound, dto.BaseResponse{Message: "参数错误"})
	}

	teacher, err := service.GetTeacherDetail(c, request.TeacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseResponse{Message: "内部错误。"})
		return
	}
	c.JSON(http.StatusOK, converter.ConvertTeacherDomainToDTO(*teacher))
}

func SearchTeacherListHandler(c *gin.Context) {}

func CreateTeacherHandler(c *gin.Context){}

func UpdateTeacherHandler(c *gin.Context){}