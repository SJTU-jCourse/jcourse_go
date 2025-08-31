package course

import "jcourse_go/internal/domain/shared"

type CourseFilter struct {
	Departments []shared.FilterItem `json:"departments"`
	Credits     []shared.FilterItem `json:"credits"`
	Semesters   []shared.FilterItem `json:"semesters"`
	Categories  []shared.FilterItem `json:"categories"`
}

type TeacherFilter struct {
	Departments []shared.FilterItem `json:"departments"`
	Titles      []shared.FilterItem `json:"titles"`
}

type TrainingPlanFilter struct {
	Departments []shared.FilterItem `json:"departments"`
	Degrees     []shared.FilterItem `json:"degrees"`
	EntryYears  []shared.FilterItem `json:"entry_years"`
}

type Pinyin struct {
	Full string `json:"full"`
	Abbr string `json:"abbr"`
}
