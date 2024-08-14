package dto

type UserRole = string

type UserType = string

type UserListRequest struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"page_size" form:"page_size"`
}

type UserListResponse = BasePaginateResponse[UserDetailDTO]

type UserListResponseForAdmin = BasePaginateResponse[UserProfileDTO]

type UserSummaryDTO struct {
	ID                   int64 `json:"id"`
	ReviewCount          int64 `json:"review_count"`
	LikeReceive          int64 `json:"like_receive"`
	TipReceive           int64 `json:"tip_receive"`
	FollowingCourseCount int64 `json:"following_course_count"`
}

type UserDetailDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

type UserProfileDTO struct {
	UserID     int64  `json:"user_id"`
	Username   string `json:"username"`
	Bio        string `json:"bio"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Role       string `json:"user_role"`
	Department string `json:"department"`
	Major      string `json:"major"`
	Grade      string `json:"grade"`
}
