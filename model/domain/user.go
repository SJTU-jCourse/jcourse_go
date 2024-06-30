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
	ID         uint
	Username   string
	Email      string
	Role       UserType // 用户在选课社区的身份
	CreatedAt  time.Time
	LastSeenAt time.Time
}

type UserProfile struct {
	UserID   uint
	Username string
	Email    string
	Avatar   string
	Type     string // 用户在学校的身份
	Major    string
	Degree   string
	Grade    string
}
