package dto

import "jcourse_go/model/model"

type CourseDetailRequest struct {
	CourseID int64 `uri:"courseID" binding:"required"`
}

type BaseCourseDetailRequest struct {
	Code string `uri:"code" binding:"required"`
}

type CourseListRequest struct {
	Page          int64  `json:"page" form:"page"`
	PageSize      int64  `json:"page_size" form:"page_size"`
	Code          string `json:"code" form:"code"`
	MainTeacherID int64  `json:"main_teacher_id" form:"main_teacher_id"`
	Departments   string `json:"departments" form:"departments"`
	Categories    string `json:"categories" form:"categories"`
	Credits       string `json:"credits" form:"credits"`
	SearchQuery   string `json:"search_query" form:"search_query"`
}

type CourseListResponse = BasePaginateResponse[model.CourseSummary]
