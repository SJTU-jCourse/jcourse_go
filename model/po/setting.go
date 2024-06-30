package po

import "gorm.io/gorm"

type SettingItemPO struct {
	gorm.Model
	Key   string
	Type  string
	Value string
}

func (po *SettingItemPO) TableName() string {
	return "setting_items"
}
