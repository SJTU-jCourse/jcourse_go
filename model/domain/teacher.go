package domain

// FIXME: Domain层的Teacher结构体中暂时还没有存Review相关的信息
type Teacher struct {
	ID         int64
	Name       string
	Email      string
	Code       string
	Department string
	Title      string
	Picture    string
	ProfileURL string
	Biography  string
	Courses    []OfferedCourse
}

type TeacherListFilter struct {
	Page             int64
	PageSize         int64
	Name             string
	Code             string
	Department       string
	Title            string
	Pinyin           string
	PinyinAbbr       string
	ContainCourseIDs []int64
}
