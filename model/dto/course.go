package dto

import "jcourse_go/model/model"

type CourseDetailRequest struct {
	CourseID int64 `uri:"courseID" binding:"required"`
}

type CourseListRequest struct {
	Page        int64  `json:"page" form:"page"`
	PageSize    int64  `json:"page_size" form:"page_size"`
	Departments string `json:"departments" form:"departments"`
	Categories  string `json:"categories" form:"categories"`
	Credits     string `json:"credits" form:"credits"`
	SearchQuery string `json:"search_query" form:"search_query"`
}

type CourseListResponse = BasePaginateResponse[model.CourseSummary]
