package shared

type UserRole string

const (
	UserRoleNormal UserRole = "normal"
	UserRoleAdmin  UserRole = "admin"
)

type IDType int64
