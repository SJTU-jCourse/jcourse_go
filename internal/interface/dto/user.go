package dto

import (
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/user"
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
	course.PaginationFilterForQuery
}

type UserListResponse = BasePaginateResponse[user.UserMinimal]
