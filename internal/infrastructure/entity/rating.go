package entity

import "time"

type RatingPO struct {
	ID int64 `gorm:"primarykey"`

	UserID      int64  `gorm:"index;index:uniq_rating,unique"`
	RelatedType string `gorm:"index;index:uniq_rating,unique;index:related_rating"`
	RelatedID   int64  `gorm:"index;index:uniq_rating,unique;index:related_rating"`
	Rating      int64  `gorm:"index"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}

func (r *RatingPO) TableName() string {
	return "ratings"
}
