package dto

import (
	"time"
)

type UserRole = string

type UserType = string

type UserListRequest struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"page_size" form:"page_size"`
}

type UserListResponse = BasePaginateResponse[UserSummaryDTO]

// 用户概要信息（用于用户列表）
type UserSummaryDTO struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Avatar   string   `json:"avatar"`
	Role     UserRole `json:"user_role"`
}

// 用户详细信息（用于用户详情页面）
type UserDetailsDTO struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Role              string    `json:"user_role"`
	LastSeenAt        time.Time `json:"lastSeenAt"`
	Type              string    `json:"user_type"`
	Avatar            string    `json:"avatar"`
	PersonalSignature string    `json:"personal_signature"`
}

// 用户个人资料信息（用于个人资料页面）
type UserProfileDTO struct {
	ID                int64  `json:"id"`
	UserID            int64  `json:"user_id"`
	Avatar            string `json:"avatar"`
	Department        string `json:"department"`
	Type              string `json:"type"`
	Major             string `json:"major"`
	Degree            string `json:"degree"`
	Grade             string `json:"grade"`
	PersonalSignature string `json:"personal_signature"`
	Username          string `json:"username"`
	Role              string `json:"user_role"`
}
