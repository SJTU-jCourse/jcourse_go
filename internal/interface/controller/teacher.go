package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/domain/course"
)

type TeacherController struct {
	teacherRepo course.TeacherRepository
}

func NewTeacherController(
	teacherRepo course.TeacherRepository,
) *TeacherController {
	return &TeacherController{
		teacherRepo: teacherRepo,
	}
}

func (c *TeacherController) GetTeacherList(ctx *gin.Context) {

}

func (c *TeacherController) GetTeacherDetail(ctx *gin.Context) {

}
