package po

import "gorm.io/gorm"

type SettingPO struct {
	gorm.Model
	Key       string `gorm:"index:uniq_setting,unique"`
	Type      string
	Value     string
	UpdatedBy int64 // user id
}

func (po *SettingPO) TableName() string {
	return "settings"
}
