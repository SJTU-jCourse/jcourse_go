package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/interface/dto"
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
	var req course.TeacherListQuery
	if err := ctx.ShouldBind(&req); err != nil {
		dto.WriteBadArgumentResponse(ctx)
		return
	}
	courses, err := c.teacherQuery.GetTeacherList(ctx, req)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, courses)
}

func (c *TeacherController) GetTeacherDetail(ctx *gin.Context) {

}

func (c *TeacherController) GetTeacherFilter(ctx *gin.Context) {
	filter, err := c.teacherQuery.GetTeacherFilter(ctx)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, filter)
}
