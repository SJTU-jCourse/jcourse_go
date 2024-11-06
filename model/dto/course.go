package dto

import "jcourse_go/model/model"

type CourseDetailRequest struct {
	CourseID int64 `uri:"courseID" binding:"required"`
}

type BaseCourseDetailRequest struct {
	Code string `uri:"code" binding:"required"`
}

type CourseListRequest struct {
	model.PaginationFilterForQuery
	Code          string `json:"code" form:"code"`
	MainTeacherID int64  `json:"main_teacher_id" form:"main_teacher_id"`
	Departments   string `json:"departments" form:"departments"`
	Categories    string `json:"categories" form:"categories"`
	Credits       string `json:"credits" form:"credits"`
}

type CourseListResponse = BasePaginateResponse[model.CourseSummary]
