package dto

import "jcourse_go/model/model"

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
	Page        int64  `json:"page" form:"page"`
	PageSize    int64  `json:"page_size" form:"page_size"`
	SearchQuery string `json:"search_query" form:"search_query"`
}

type UserListResponse = BasePaginateResponse[model.UserMinimal]
