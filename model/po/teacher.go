package po

import "gorm.io/gorm"

type TeacherPO struct {
	gorm.Model
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
}

func (po *TeacherPO) TableName() string {
	return "teachers"
}
