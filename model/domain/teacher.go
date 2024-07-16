package domain

type Teacher struct {
	ID         int64
	Name       string
	Email      string
	Code       string
	Department string
	Title      string
	Courses    []OfferedCourse
}

type TeacherListFilter struct{
	Name string
	Department string
	Pinyin string
	PinyinAbbr string
}