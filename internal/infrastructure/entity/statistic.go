package entity

import "time"

type Statistic struct {
	ID int64 `gorm:"primaryKey"`

	Date            string `gorm:"index:unique"` // 异步落盘, yyyy-mm-dd
	DailyNewUser    int64
	DailyNewReview  int64
	DailyActiveUser int64
	DailyPageView   int64
	TotalUser       int64
	TotalReview     int64

	CreatedAt time.Time
}

func (po *Statistic) TableName() string {
	return "statistic"
}
