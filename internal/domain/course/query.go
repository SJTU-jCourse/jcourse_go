package course

import "jcourse_go/internal/domain/shared"

type CourseListQuery struct {
	shared.PaginationQuery
	Code          string
	MainTeacherID int64
	Departments   []string
	Categories    []string
	Credits       []float64
	Semesters     []string
}

type TeacherListQuery struct {
	shared.PaginationQuery
	Name             string
	Code             string
	Departments      []string
	Titles           []string
	Pinyin           string
	PinyinAbbr       string
	ContainCourseIDs []int64
}

type TrainingPlanListQuery struct {
	shared.PaginationQuery
	Major            string
	Degrees          []string
	Departments      []string
	EntryYears       []string
	ContainCourseIDs []int64
}

type ReviewQuery struct {
	shared.PaginationQuery
	CourseID         int64
	Semester         string
	UserID           int64
	ReviewID         int64
	Rating           int64
	ExcludeAnonymous bool
}
