package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/interface/dto"
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
	var req dto.CourseDetailRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.WriteBadArgumentResponse(ctx)
		return
	}

	course, err := c.courseQuery.GetCourseDetail(ctx, shared.IDType(req.CourseID))
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}

	dto.WriteDataResponse(ctx, course)
}

func (c *CourseController) GetCourseList(ctx *gin.Context) {
	var req course.CourseListQuery
	if err := ctx.ShouldBind(&req); err != nil {
		dto.WriteBadArgumentResponse(ctx)
		return
	}
	courses, err := c.courseQuery.GetCourseList(ctx, req)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, courses)
}

func (c *CourseController) GetCourseFilter(ctx *gin.Context) {
	filter, err := c.courseQuery.GetCourseFilter(ctx)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, filter)
}

func (c *CourseController) SubscribeCourse(ctx *gin.Context) {

}
