package dto

import "jcourse_go/model/domain"

type CourseListDTO struct {
	ID              int64    `json:"id"`
	Code            string   `json:"code"`
	Name            string   `json:"name"`
	Credit          float64  `json:"credit"`
	MainTeacherName string   `json:"main_teacher_name"`
	Categories      []string `json:"categories"`
	Department      string   `json:"department"`
}

type CourseListRequest = domain.CourseListFilter

type CourseListResponse = BasePaginateResponse[CourseListDTO]
