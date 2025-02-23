package po

import (
	"time"
)

type UserPointDetailPO struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`

	EventType   string  `gorm:"index"`
	Description string  `gorm:"index"`
	Value       int64   // 积分变动值
	UserID      int64   `gorm:"index"` // 用户ID
	User        *UserPO `gorm:"constraint:OnDelete:CASCADE;"`
}

func (po *UserPointDetailPO) TableName() string { return "user_point_details" }
