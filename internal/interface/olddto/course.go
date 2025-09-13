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

type CourseListResponse = BasePaginateResponse[course.CourseSummary]
