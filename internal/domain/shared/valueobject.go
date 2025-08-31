package shared

type FilterItem struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

type UserRole string

const (
	UserRoleNormal UserRole = "normal"
	UserRoleAdmin  UserRole = "admin"
)
