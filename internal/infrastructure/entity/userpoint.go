package entity

import "time"

type UserPoint struct {
	ID int64 `gorm:"primaryKey"`

	UserID int64 `gorm:"index"` // 用户ID
	User   *User

	Type        string `gorm:"index"`
	Description string
	Value       int64 // 积分变动值

	CreatedAt time.Time `gorm:"index"`
}

func (po *UserPoint) TableName() string { return "user_point" }
