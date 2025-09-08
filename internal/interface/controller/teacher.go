package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application"
)

type TeacherController struct {
	teacherQuery application.TeacherQueryService
}

func NewTeacherController(
	teacherQuery application.TeacherQueryService,
) *TeacherController {
	return &TeacherController{
		teacherQuery: teacherQuery,
	}
}

func (c *TeacherController) GetTeacherList(ctx *gin.Context) {

}

func (c *TeacherController) GetTeacherDetail(ctx *gin.Context) {

}
