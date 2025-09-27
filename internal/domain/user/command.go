package user

import "jcourse_go/internal/domain/shared"

type UserQuery struct {
	shared.PaginationQuery
}

type UserPointQuery struct {
	shared.PaginationQuery
	UserPointDetailID int64
	UserID            int64
	StartTime         int64
	EndTime           int64
}

type UpdateUserInfoCommand struct {
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}
