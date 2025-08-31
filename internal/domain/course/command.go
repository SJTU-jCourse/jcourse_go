package course

import "jcourse_go/internal/domain/shared"

type CourseListFilterForQuery struct {
	shared.PaginationFilterForQuery
	Code          string
	MainTeacherID int64
	Departments   []string
	Categories    []string
	Credits       []float64
}

type TeacherFilterForQuery struct {
	shared.PaginationFilterForQuery
	Name             string
	Code             string
	Departments      []string
	Titles           []string
	Pinyin           string
	PinyinAbbr       string
	ContainCourseIDs []int64
}

type TrainingPlanFilterForQuery struct {
	shared.PaginationFilterForQuery
	Major            string
	Degrees          []string
	Departments      []string
	EntryYears       []string
	ContainCourseIDs []int64
}

type ReviewFilterForQuery struct {
	shared.PaginationFilterForQuery
	CourseID         int64
	Semester         string
	UserID           int64
	ReviewID         int64
	Rating           int64
	ExcludeAnonymous bool
}
