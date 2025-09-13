package olddto

import (
	"jcourse_go/internal/domain/course"
)

type TeacherDetailRequest struct {
	TeacherID int64 `uri:"teacherID" binding:"required"`
}

type TeacherListRequest struct {
	course.PaginationFilterForQuery
	Name        string `json:"name" form:"name"`
	Code        string `json:"code" form:"code"`
	Departments string `json:"departments" form:"departments"`
	Titles      string `json:"titles" form:"titles"`
	Pinyin      string `json:"pinyin" form:"pinyin"`
	PinyinAbbr  string `json:"pinyin_abbr" form:"pinyin_abbr"`
}

type TeacherListResponse = BasePaginateResponse[course.TeacherSummary]
