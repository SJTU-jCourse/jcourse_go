package model

type UserRole = string

const (
	UserRoleNormal UserRole = "normal"
	UserRoleAdmin  UserRole = "admin"
)

type UserType = string

const (
	UserTypeStudent UserType = "student"
	UserTypeFaculty UserType = "faculty"
)

type UserFilterForQuery struct {
	PaginationFilterForQuery
}

type UserDetail struct {
	UserMinimal
	Bio        string `json:"bio"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Type       string `json:"type"`
	Department string `json:"department"`
	Major      string `json:"major"`
	Grade      string `json:"grade"`
	Points     int64  `json:"points"`
}

type UserActivity struct {
	ReviewCount          int64 `json:"review_count"`
	LikeReceive          int64 `json:"like_receive"`
	TipReceive           int64 `json:"tip_receive"`
	FollowingCourseCount int64 `json:"following_course_count"`
}

type UserMinimal struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
