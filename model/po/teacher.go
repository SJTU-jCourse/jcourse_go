package po

import "gorm.io/gorm"

type TeacherPO struct {
	gorm.Model
	Name       string `gorm:"index"`
	Code       string `gorm:"index;uniqueIndex"`
	Email      string `gorm:"index"`
	Department string `gorm:"index"`
	Title      string
	Pinyin     string `gorm:"index"`
	PinyinAbbr string `gorm:"index"`
}
