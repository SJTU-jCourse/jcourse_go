package domain

type Teacher struct {
	ID         int64
	Name       string
	Email      string
	Code       string
	Department string
	Title      string
	Picture    string
	ProfileURL string
	Courses    []BaseCourse
}

type TeacherFilter struct {
	Name             string
	Code             string
	Department       string
	Title            string
	ContainCourseIDs []int64
}

type TeacherListFilter struct{
	Name string
	Department string
	Pinyin string
	PinyinAbbr string
}