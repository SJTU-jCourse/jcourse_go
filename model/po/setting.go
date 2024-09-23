package po

import "gorm.io/gorm"

type SettingPO struct {
	gorm.Model
	Key       string `gorm:"index:uniq_setting,unique"`
	Type      string
	Value     string
	UpdatedBy int64 // user id
	Client    bool  // should client side fetch
}

func (po *SettingPO) TableName() string {
	return "settings"
}
