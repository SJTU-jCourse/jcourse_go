package entity

import "time"

type Statistic struct {
	ID int64 `gorm:"primaryKey"`

	Date         string `gorm:"index:unique"` // 异步落盘, yyyy-mm-dd
	NewUsers     int64
	NewReviews   int64
	ActiveUsers  int64
	TotalUsers   int64
	TotalReviews int64

	CreatedAt time.Time
}

func (po *Statistic) TableName() string {
	return "statistic"
}
