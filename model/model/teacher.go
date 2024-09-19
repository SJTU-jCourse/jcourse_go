package model

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
	SearchQuery      string
}

type TeacherDetail struct {
	TeacherSummary
	Email      string          `json:"email"`
	Code       string          `json:"code"`
	Title      string          `json:"title"`
	ProfileURL string          `json:"profile_url"`
	Biography  string          `json:"biography"`
	Courses    []OfferedCourse `json:"courses"`
}

type TeacherSummary struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Department string     `json:"department"`
	Picture    string     `json:"picture"`
	ReviewInfo RatingInfo `json:"rating_info"`
}
