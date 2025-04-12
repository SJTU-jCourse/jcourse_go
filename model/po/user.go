package po

import (
	"time"
)

type UserPO struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`

	Username string `gorm:"index:idx_auth;uniqueIndex"`
	Email    string `gorm:"uniqueIndex"`
	Password string `gorm:"index:idx_auth"`
	UserRole string `gorm:"index"` // 用户在选课社区的身份

	Type string // 用户在学校的身份
	// Department string // 院系
	// Major      string // 专业
	// Degree     string // 学位
	// Grade      string // 年级

	Avatar string // 头像
	Bio    string // 个人介绍

	Points           int64                // 积分
	UserPointDetails []*UserPointDetailPO `gorm:"foreignkey:UserID"`

	LastSeenAt time.Time
}

func (po *UserPO) TableName() string {
	return "users"
}
