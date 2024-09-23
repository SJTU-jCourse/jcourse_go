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

type CourseListFilter struct {
	Page          int64
	PageSize      int64
	Code          string
	MainTeacherID int64
	Departments   []string
	Categories    []string
	Credits       []float64
	SearchQuery   string
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
