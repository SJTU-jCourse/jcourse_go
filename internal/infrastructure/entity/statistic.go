package entity

import "time"

type StatisticPO struct {
	ID int64 `gorm:"primarykey"`

	Date         string `gorm:"index:unique"` // 异步落盘, yyyy-mm-dd
	NewUsers     int64
	NewReviews   int64
	ActiveUsers  int64
	TotalUsers   int64
	TotalReviews int64

	CreatedAt time.Time `gorm:"index"`
}

func (po *StatisticPO) TableName() string {
	return "statistics"
}
