package entity

import "time"

type TeacherPO struct {
	ID int64 `gorm:"primarykey"`

	Name       string `gorm:"index"`
	Pinyin     string `gorm:"index"`
	PinyinAbbr string `gorm:"index"`

	Code       string `gorm:"index:,unique"`
	Email      string `gorm:"index"`
	Department string `gorm:"index"`
	Title      string

	Picture    string // picture URL
	ProfileURL string
	Biography  string // 个人简述

	SearchIndex SearchIndex `gorm:"->:false;<-"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
}

func (po *TeacherPO) TableName() string {
	return "teachers"
}
