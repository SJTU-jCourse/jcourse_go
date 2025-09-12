package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
)

type CourseController struct {
	courseQuery query.CourseQueryService
}

func NewCourseController(
	courseQuery query.CourseQueryService,
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
