package user

import "jcourse_go/internal/domain/shared"

type User struct {
	ID       int64           `json:"id"`
	Username string          `json:"username"`
	Bio      string          `json:"bio"`
	Email    string          `json:"email"`
	Role     shared.UserRole `json:"role"`
}

// UserPoint 用户积分明细
type UserPoint struct {
	ID          int64         `json:"id"`
	Value       int64         `json:"value"` // 积分变化值: +1, -3
	Type        UserPointType `json:"type"`
	Description string        `json:"description"`
	CreatedAt   string        `json:"created_at"`
}
