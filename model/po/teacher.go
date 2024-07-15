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
	Picture    string `gorm:"index"`
	ProfileURL string `gorm:"index"`
}

func (po *TeacherPO) TableName() string {
	return "teachers"
}

type TeacherCoursePO struct {
	gorm.Model
	TeacherID int64 `gorm:"index"`
	CourseID  int64 `gorm:"index"`
}
