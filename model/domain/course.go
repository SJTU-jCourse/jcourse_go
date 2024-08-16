package domain

type OfferedCourse struct {
	ID           int64
	TeacherGroup []Teacher
	Semester     string
	Language     string
	Grade        []string
}

type RatingInfoDistItem struct {
	Rating int64
	Count  int64
}

type RatingInfo struct {
	Average    float64              `json:"average"`
	Count      int64                `json:"count"`
	RatingDist []RatingInfoDistItem `json:"rating_dist"`
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
	RatingInfo     RatingInfo
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
	MajorCode  string
	MajorClass string
	Department string
	EntryYear  string
	MinCredits float64
	TotalYear  int
	Courses    []TrainingPlanCourse
}
type TrainingPlanCourse struct {
	ID              int64
	Code            string
	Name            string
	Credit          float64
	SuggestSemester string
	Department      string
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
