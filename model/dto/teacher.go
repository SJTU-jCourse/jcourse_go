package dto

import "jcourse_go/model/model"

type TeacherDetailRequest struct {
	TeacherID int64 `uri:"teacherID" binding:"required"`
}

type TeacherListRequest struct {
	model.PaginationFilterForQuery
	Name        string `json:"name" form:"name"`
	Code        string `json:"code" form:"code"`
	Departments string `json:"departments" form:"departments"`
	Titles      string `json:"titles" form:"titles"`
	Pinyin      string `json:"pinyin" form:"pinyin"`
	PinyinAbbr  string `json:"pinyin_abbr" form:"pinyin_abbr"`
}

type TeacherListResponse = BasePaginateResponse[model.TeacherSummary]
