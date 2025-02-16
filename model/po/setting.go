package po

import (
	"time"
)

type SettingPO struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`

	Key       string `gorm:"index:uniq_setting,unique"`
	Type      string
	Value     string
	UpdatedBy int64 `gorm:"index"` // user id
	Client    bool  // should client side fetch
}

func (po *SettingPO) TableName() string {
	return "settings"
}
