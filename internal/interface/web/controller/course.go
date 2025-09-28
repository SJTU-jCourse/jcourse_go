package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/command"
	"jcourse_go/internal/application/query"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/notification"
	"jcourse_go/internal/domain/shared"
)

type CourseController struct {
	courseQuery        query.CourseQueryService
	courseNotification command.CourseNotificationService
}

func NewCourseController(
	courseQuery query.CourseQueryService,
	courseNotification command.CourseNotificationService,
) *CourseController {
	return &CourseController{
		courseQuery:        courseQuery,
		courseNotification: courseNotification,
	}
}

func (c *CourseController) GetCourseDetail(ctx *gin.Context) {
	courseIDStr := ctx.Param("courseID")
	if courseIDStr == "" {
		WriteBadArgumentResponse(ctx)
		return
	}
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil || courseID == 0 {
		WriteBadArgumentResponse(ctx)
		return
	}

	course, err := c.courseQuery.GetCourseDetail(ctx, shared.IDType(courseID))
	if err != nil {
		WriteErrorResponse(ctx, err)
		return
	}

	WriteDataResponse(ctx, course)
}

func (c *CourseController) GetCourseList(ctx *gin.Context) {
	var req course.CourseListQuery
	if err := ctx.ShouldBind(&req); err != nil {
		WriteBadArgumentResponse(ctx)
		return
	}
	courses, err := c.courseQuery.GetCourseList(ctx, req)
	if err != nil {
		WriteErrorResponse(ctx, err)
		return
	}
	WriteDataResponse(ctx, courses)
}

func (c *CourseController) GetCourseFilter(ctx *gin.Context) {
	filter, err := c.courseQuery.GetCourseFilter(ctx)
	if err != nil {
		WriteErrorResponse(ctx, err)
		return
	}
	WriteDataResponse(ctx, filter)
}

func (c *CourseController) ChangeNotification(ctx *gin.Context) {
	var req notification.CourseNotificationCommand
	if err := ctx.ShouldBind(&req); err != nil {
		WriteErrorResponse(ctx, err)
		return
	}
	reqCtx := shared.NewRequestCtx(0, shared.UserRoleNormal)
	if err := c.courseNotification.Change(ctx, reqCtx, req); err != nil {
		WriteErrorResponse(ctx, err)
		return
	}
	WriteDataResponse(ctx, nil)
}
