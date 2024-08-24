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

	SearchIndex SearchIndex `gorm:"index:idx_search, type:gin"`
}

func (po *UserPO) TableName() string {
	return "users"
}

type UserProfilePO struct {
	gorm.Model
	UserID     int64
	Avatar     string
	Department string
	Type       string // 用户在学校的身份
	Major      string
	Degree     string
	Grade      string
	Bio        string
}

func (profile *UserProfilePO) TableName() string {
	return "user_profiles"
}

type UserActivityPO struct {
	gorm.Model
	UserID       int64     // 用户ID
	ActivityType string    // 活动类型，如发布课程点评、点赞、回复、关注/屏蔽用户/课程等。
	TargetID     string    // 活动对象ID
	CreatedAt    time.Time // 活动发生时间
}

func (userActivity *UserActivityPO) TableName() string { return "user_activities" }
