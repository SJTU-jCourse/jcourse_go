package course

import "jcourse_go/internal/domain/shared"

type CourseListQuery struct {
	shared.PaginationQuery
	Code          string    `json:"code" form:"code"`
	MainTeacherID int64     `json:"main_teacher_id" form:"main_teacher_id"`
	Departments   []string  `json:"departments" form:"departments"`
	Categories    []string  `json:"categories" form:"categories"`
	Credits       []float64 `json:"credits" form:"credits"`
	Semesters     []string  `json:"semesters" form:"semesters"`
}

type TeacherListQuery struct {
	shared.PaginationQuery
	Name        string   `json:"name" form:"name"`
	Code        string   `json:"code" form:"code"`
	Departments []string `json:"departments" form:"departments"`
	Titles      []string `json:"titles" form:"titles"`
}

type TrainingPlanListQuery struct {
	shared.PaginationQuery
	Major       string   `json:"major" form:"major"`
	Degrees     []string `json:"degrees" form:"degrees"`
	Departments []string `json:"departments" form:"departments"`
	EntryYears  []string `json:"entry_years" form:"entry_years"`
}
