package po

import (
	"time"

	"gorm.io/gorm"
)

type StatisticPO struct {
	gorm.Model
	UVCount      int64
	PVCount      int64
	Date         time.Time `gorm:"index:,unique"` // createdAt time != actual datetime, 可能异步落盘, 统一为中午12点
	NewCourses   int64
	NewUsers     int64
	NewReviews   int64
	TotalUsers   int64 // 这几个是在db之中存储,按天更新,还是调用接口实时统计?
	TotalReviews int64
	TotalCourses int64
}

// StatisticDataPO uv bitmap, etc, separate date and count
type StatisticDataPO struct {
	gorm.Model
	StatisticID int64     `gorm:"index:,unique"`
	Date        time.Time `gorm:"index:,unique"` // redundant column for speed up
	UVData      []byte
}

func (po *StatisticPO) TableName() string {
	return "statistics"
}
func (po *StatisticDataPO) TableName() string { return "statistics_data" }
