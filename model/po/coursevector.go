package po

import "github.com/lib/pq"

type CourseVectorPO struct {
	ID     int64           `gorm:"primaryKey"`
	Vector pq.Float32Array `gorm:"type:float4[]"`
}

func (po *CourseVectorPO) TableName() string {
	return "course_vectors"
}
