package domain

type OfferedCourse struct {
	ID           int64
	Course       *Course
	TeacherGroup []Teacher
	Semester     string
	Language     string
	Grade        []string
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
}

type TrainingPlan struct {
	ID         int64
	Major      string
	Department string
	EntryYear  string
	Courses    []BaseCourse
}

type BaseCourse struct {
	ID     int64
	Code   string
	Name   string
	Credit float64
}

type CourseListFilter struct {
	Page        int64     `json:"page" form:"page"`
	PageSize    int64     `json:"page_size" form:"page_size"`
	Departments []string  `json:"departments" form:"departments"`
	Categories  []string  `json:"categories" form:"categories"`
	Credits     []float64 `json:"credits" form:"credits"`
}
