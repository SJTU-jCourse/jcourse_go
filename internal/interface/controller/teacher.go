package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
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
	teacherIDStr := ctx.Param("teacherID")
	if teacherIDStr == "" {
		dto.WriteBadArgumentResponse(ctx)
		return
	}
	teacherID, err := strconv.Atoi(teacherIDStr)
	if err != nil || teacherID == 0 {
		dto.WriteBadArgumentResponse(ctx)
		return
	}

	teacher, err := c.teacherQuery.GetTeacherDetail(ctx, shared.IDType(teacherID))
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, teacher)
}

func (c *TeacherController) GetTeacherFilter(ctx *gin.Context) {
	filter, err := c.teacherQuery.GetTeacherFilter(ctx)
	if err != nil {
		dto.WriteErrorResponse(ctx, err)
		return
	}
	dto.WriteDataResponse(ctx, filter)
}
