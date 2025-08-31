package shared

type UserRole string

const (
	UserRoleNormal UserRole = "normal"
	UserRoleAdmin  UserRole = "admin"
)

type IDType int64

func (i IDType) Int64() int64 {
	return int64(i)
}
