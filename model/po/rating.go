package po

import (
	"time"
)

type RatingPO struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time

	UserID      int64  `gorm:"index;index:uniq_rating,unique"`
	User        UserPO `gorm:"constraint:OnDelete:CASCADE;"`
	RelatedType string `gorm:"index;index:uniq_rating,unique"`
	RelatedID   int64  `gorm:"index;index:uniq_rating,unique"`
	Rating      int64  `gorm:"index"`
}

func (r *RatingPO) TableName() string {
	return "ratings"
}

type RatingInfo struct {
	Average float64
	Count   int64
}
