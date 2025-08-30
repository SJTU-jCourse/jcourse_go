package dto

import (
	model2 "jcourse_go/internal/model/model"
)

type UserProfileDTO struct {
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	Bio        string `json:"bio"`
	Type       string `json:"type"`
	Department string `json:"department"`
	Major      string `json:"major"`
	Grade      string `json:"grade"`
	Degree     string `json:"degree"`
}

type UserListRequest struct {
	model2.PaginationFilterForQuery
}

type UserListResponse = BasePaginateResponse[model2.UserMinimal]
