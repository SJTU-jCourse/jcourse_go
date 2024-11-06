package model

type TeacherFilterForQuery struct {
	PaginationFilterForQuery
	Name             string
	Code             string
	Departments      []string
	Titles           []string
	Pinyin           string
	PinyinAbbr       string
	ContainCourseIDs []int64
}

type TeacherDetail struct {
	TeacherSummary
	Email      string          `json:"email"`
	Code       string          `json:"code"`
	Title      string          `json:"title"`
	ProfileURL string          `json:"profile_url"`
	Biography  string          `json:"biography"`
	Courses    []CourseSummary `json:"courses"`
}

type TeacherSummary struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Department string     `json:"department"`
	Picture    string     `json:"picture"`
	RatingInfo RatingInfo `json:"rating_info"`
}

type TeacherFilter struct {
	Departments []FilterItem `json:"departments"`
	Titles      []FilterItem `json:"titles"`
}
