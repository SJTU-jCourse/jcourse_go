package entity

import (
	"time"
)

type User struct {
	ID int64 `gorm:"primaryKey"`

	Username string `gorm:"index:idx_auth;uniqueIndex"`
	Email    string `gorm:"index:idx_email;uniqueIndex"`
	Password string `gorm:"index:idx_auth"`
	UserRole string `gorm:"index"` // 用户在选课社区的身份

	LowerCase bool // 历史逻辑

	LastSeenAt    time.Time
	SuspendedTill *time.Time

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
}

func (po *User) TableName() string {
	return "user"
}
