package user

import (
	"time"

	"jcourse_go/internal/domain/shared"
)

type User struct {
	ID       shared.IDType   `json:"id"`
	Username string          `json:"username"`
	Bio      string          `json:"bio"`
	Email    string          `json:"email"`
	Role     shared.UserRole `json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserPoint 用户积分明细
type UserPoint struct {
	ID          shared.IDType  `json:"id"`
	Value       int64          `json:"value"` // 积分变化值: +1, -3
	Event       PointEventType `json:"event"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
}

type UserPointTransaction struct {
	FromID    shared.IDType `json:"from_id"`
	ToID      shared.IDType `json:"to_id"`
	Value     int64         `json:"value"`
	Reason    string        `json:"reason"`
	CreatedAt time.Time     `json:"created_at"`
}
