package entity

import "time"

type ApiKey struct {
	ID          int64  `gorm:"primaryKey"`
	Key         string `gorm:"uniqueIndex"`
	UserID      int64  `gorm:"index:uniq_user,unique"`
	Description string
	CreatedAt   time.Time
	DeletedAt   *time.Time
}

func (ApiKey) TableName() string {
	return "api_key"
}
