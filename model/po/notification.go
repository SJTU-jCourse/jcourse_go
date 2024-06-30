package po

import "gorm.io/gorm"

type NotificationPO struct {
	gorm.Model
	UserID      int64
	RelatedType string
	RelatedID   int64
	Title       string
	Content     string
	Read        bool `gorm:"index"`
}

func (po *NotificationPO) TableName() string {
	return "notifications"
}
