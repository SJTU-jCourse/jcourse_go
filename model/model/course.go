package model

type BaseCourse struct {
	ID     int64   `json:"id"`
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Credit float64 `json:"credit"`
}

type OfferedCourse struct {
	ID           int64           `json:"id"`
	Semester     string          `json:"semester"`
	Grade        []string        `json:"grade"`
	Language     string          `json:"language"`
	TeacherGroup []TeacherDetail `json:"teacher_group"`
}

type CourseListFilterForQuery struct {
	PaginationFilterForQuery
	Code          string
	MainTeacherID int64
	Departments   []string
	Categories    []string
	Credits       []float64
}

type CourseMinimal struct {
	BaseCourse
	ID          int64          `json:"id"`
	MainTeacher TeacherSummary `json:"main_teacher"`
}

type CourseSummary struct {
	CourseMinimal
	Categories []string   `json:"categories"`
	Department string     `json:"department"`
	RatingInfo RatingInfo `json:"rating_info"`
}

type CourseDetail struct {
	CourseSummary
	OfferedCourse []OfferedCourse `json:"offered_courses"`
}

type FilterItem struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

type CourseFilter struct {
	Departments []FilterItem `json:"departments"`
	Credits     []FilterItem `json:"credits"`
	Semesters   []FilterItem `json:"semesters"`
	Categories  []FilterItem `json:"categories"`
}

type PaginationFilterForQuery struct {
	Page      int64  `json:"page" form:"page"`
	PageSize  int64  `json:"page_size" form:"page_size"`
	Search    string `json:"search" form:"search"`
	Order     string `json:"order" form:"order"`
	Ascending bool   `json:"ascending" form:"ascending"`
}
