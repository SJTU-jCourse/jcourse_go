package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/domain/course"
)

type CourseController struct {
	courseRepo course.CourseRepository
}

func NewCourseController(
	courseRepo course.CourseRepository,
) *CourseController {
	return &CourseController{
		courseRepo: courseRepo,
	}
}

func (c *CourseController) GetCourseDetail(ctx *gin.Context) {

}

func (c *CourseController) GetCourseList(ctx *gin.Context) {

}

func (c *CourseController) GetCourseFilter(ctx *gin.Context) {

}

func (c *CourseController) SubscribeCourse(ctx *gin.Context) {

}
