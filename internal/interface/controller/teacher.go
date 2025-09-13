package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
)

type TeacherController struct {
	teacherQuery query.TeacherQueryService
}

func NewTeacherController(
	teacherQuery query.TeacherQueryService,
) *TeacherController {
	return &TeacherController{
		teacherQuery: teacherQuery,
	}
}

func (c *TeacherController) GetTeacherList(ctx *gin.Context) {

}

func (c *TeacherController) GetTeacherDetail(ctx *gin.Context) {

}

func (c *TeacherController) GetTeacherFilter(ctx *gin.Context) {}
