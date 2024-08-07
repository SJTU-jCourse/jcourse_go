package domain

type OfferedCourse struct {
	ID           int64
	TeacherGroup []Teacher
	Semester     string
	Language     string
	Grade        []string
}

type CourseReviewInfo struct {
	Average float64 `json:"average"`
	Count   int64   `json:"count"`
}

type Course struct {
	ID             int64
	Code           string
	Name           string
	Credit         float64
	MainTeacher    Teacher
	Department     string
	Categories     []string
	OfferedCourses []OfferedCourse
	ReviewInfo     CourseReviewInfo
}

type TrainingPlan struct {
	ID         int64
	Major      string
	Department string
	EntryYear  string
	Courses    []BaseCourse
}
type TrainingPlanDetail struct {
	ID         int64
	Major      string
	Department string
	EntryYear  string
	MajorCode  string
	MajorClass string
	MinPoints  float64
	TotalYear  int
	Courses    []TrainingPlanCourse
}
type TrainingPlanCourse struct {
	ID              int64
	Code            string
	Name            string
	Credit          float64
	SuggestYear     int64
	SuggestSemester int64
	Department      string
}
type TrainingPlanRateInfo struct {
	Avg            float64
	Count          int64
	TrainingPlanID int64
	Rates          []TrainingPlanRate
}
type TrainingPlanRate struct {
	TrainingPlanID int64
	UserID         int64
	Rate           int64
}

type TrainingPlanFilter struct {
	Page             int64
	PageSize         int64
	Major            string
	Department       string
	EntryYear        string
	ContainCourseIDs []int64
}

type BaseCourse struct {
	ID     int64
	Code   string
	Name   string
	Credit float64
}

type CourseListFilter struct {
	Page        int64
	PageSize    int64
	Departments []string
	Categories  []string
	Credits     []float64
}
