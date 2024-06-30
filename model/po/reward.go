package po

import "gorm.io/gorm"

type RewardRecordPO struct {
	gorm.Model
	UserID       int64
	RewardType   string
	RewardAmount int64
}

func (po *RewardRecordPO) TableName() string {
	return "reward_records"
}
