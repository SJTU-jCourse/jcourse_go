package dto

import "jcourse_go/model/model"

type TeacherDetailRequest struct {
	TeacherID int64 `uri:"teacherID" binding:"required"`
}

type TeacherListRequest struct {
	Page        int64  `json:"page" form:"page"`
	PageSize    int64  `json:"page_size" form:"page_size"`
	Name        string `json:"name" form:"name"`
	Code        string `json:"code" form:"code"`
	Department  string `json:"departments" form:"departments"`
	Title       string `json:"title" form:"title"`
	Pinyin      string `json:"pinyin" form:"pinyin"`
	PinyinAbbr  string `json:"pinyin_abbr" form:"pinyin_abbr"`
	SearchQuery string `json:"search_query" form:"search_query"`
}

type TeacherListResponse = BasePaginateResponse[model.TeacherSummary]
