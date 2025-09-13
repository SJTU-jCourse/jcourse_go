package olddto

import (
	"jcourse_go/internal/domain/course"
)

type CourseDetailRequest struct {
	CourseID int64 `uri:"courseID" binding:"required"`
}

type BaseCourseDetailRequest struct {
	Code string `uri:"code" binding:"required"`
}

type CourseListRequest struct {
	course.PaginationFilterForQuery
	Code          string `json:"code" form:"code"`
	MainTeacherID int64  `json:"main_teacher_id" form:"main_teacher_id"`
	Departments   string `json:"departments" form:"departments"`
	Categories    string `json:"categories" form:"categories"`
	Credits       string `json:"credits" form:"credits"`
}

type CourseListResponse = BasePaginateResponse[course.CourseSummary]
