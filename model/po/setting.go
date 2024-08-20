package po

import "gorm.io/gorm"

type SettingItemPO struct {
	gorm.Model
	Key       string `gorm:"index"`
	Value     string
	UpdatedBy int64 // user id
}

func (po *SettingItemPO) TableName() string {
	return "settings"
}
