package po

import "gorm.io/gorm"

type UserPointDetailPO struct {
	gorm.Model
	EventType   string `gorm:"index"`
	Description string `gorm:"index"`
	Value       int64  // 积分变动值
	UserID      int64  `gorm:"index"` // 用户ID
}

func (po *UserPointDetailPO) TableName() string { return "user_point_details" }
