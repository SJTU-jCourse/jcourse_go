package po

import (
	"time"
)

type StatisticPO struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`

	UVCount      int64
	PVCount      int64
	Date         string `gorm:"index:,unique"` // 异步落盘, yyyy-mm-dd
	NewUsers     int64
	NewReviews   int64
	TotalUsers   int64 // 这几个是在db之中存储,按天更新,还是调用接口实时统计?
	TotalReviews int64
}

// StatisticDataPO uv bitmap, etc, separate date and count
type StatisticDataPO struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`

	StatisticID int64  `gorm:"index:,unique"`
	Date        string `gorm:"index:,unique"` // redundant column for speed up
	UVData      []byte
}

func (po *StatisticPO) TableName() string {
	return "statistics"
}
func (po *StatisticDataPO) TableName() string { return "statistics_data" }
