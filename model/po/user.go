package po

import (
	"time"

	"gorm.io/gorm"
)

type UserPO struct {
	gorm.Model
	Username   string `gorm:"index:idx_auth;uniqueIndex"`
	Email      string `gorm:"uniqueIndex"`
	Password   string `gorm:"index:idx_auth"`
	UserRole   string `gorm:"index"` // 用户在选课社区的身份
	LastSeenAt time.Time
}

func (po *UserPO) TableName() string {
	return "users"
}

type UserProfile struct {
	gorm.Model
	UserID int64
	Type   string // 用户在学校的身份
	Major  string
	Degree string
	Grade  string
}

func (profile *UserProfile) TableName() string {
	return "user_profiles"
}
