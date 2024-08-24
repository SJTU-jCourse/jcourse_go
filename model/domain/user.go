package domain

import "time"

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

type User struct {
	ID         int64
	Username   string
	Email      string
	Role       UserRole // 用户在选课社区的身份
	CreatedAt  time.Time
	LastSeenAt time.Time

	Profile UserProfile
}

type UserProfile struct {
	UserID     int64
	Avatar     string
	Department string
	Type       UserType // 用户在学校的身份
	Major      string
	Degree     string
	Grade      string
	Bio        string
}

type UserFilter struct {
	Page        int64
	PageSize    int64
	SearchQuery string
	// To be continued ... (add more fields)
}
