package po

import (
	"time"
)

type TeacherPO struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`

	Name       string `gorm:"index"`
	Code       string `gorm:"index:,unique"`
	Email      string `gorm:"index"`
	Department string `gorm:"index"`
	Title      string
	Pinyin     string `gorm:"index"`
	PinyinAbbr string `gorm:"index"`
	Picture    string // picture URL
	ProfileURL string
	Biography  string // 个人简述

	RatingCount int64   `gorm:"index;default:0;not null"`
	RatingAvg   float64 `gorm:"index;default:0;not null"`

	Courses []CoursePO `gorm:"foreignKey:MainTeacherID"`

	SearchIndex SearchIndex `gorm:"->:false;<-"`
}

func (po *TeacherPO) TableName() string {
	return "teachers"
}
