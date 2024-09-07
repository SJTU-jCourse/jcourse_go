package dto

type TeacherDTO struct {
	ID         int64              `json:"id"`
	Email      string             `json:"email"`
	Code       string             `json:"code"`
	Name       string             `json:"name"`
	Department string             `json:"department"`
	Title      string             `json:"title"`
	Picture    string             `json:"picture"`
	ProfileURL string             `json:"profile_url"`
	Biography  string             `json:"biography"`
	Courses    []OfferedCourseDTO `json:"courses"`
}

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

type TeacherListResponse = BasePaginateResponse[TeacherDTO]
