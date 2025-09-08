package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application"
)

type CourseController struct {
	courseQuery application.CourseQueryService
}

func NewCourseController(
	courseQuery application.CourseQueryService,
) *CourseController {
	return &CourseController{
		courseQuery: courseQuery,
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
