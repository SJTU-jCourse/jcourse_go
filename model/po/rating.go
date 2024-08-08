package po

import "gorm.io/gorm"

type RatingPO struct {
	gorm.Model
	UserID      int64  `gorm:"index;index:uniq_rating,unique"`
	RelatedType string `gorm:"index;index:uniq_rating,unique"`
	RelatedID   int64  `gorm:"index;index:uniq_rating,unique"`
	Rating      int64  `gorm:"index"`
}

func (r *RatingPO) TableName() string {
	return "ratings"
}
