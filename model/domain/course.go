package domain

type CourseLike interface {
	GetCode() string
	GetName() string
	GetCredit() float32
}

type OfferedCourse struct {
	ID int64
	Course
	MainTeacher  Teacher
	TeacherGroup []Teacher
	Semester     string
	Language     string
	Grade        []string
}

type Course struct {
	ID          int64
	Code        string
	Name        string
	Credit      float64
	MainTeacher Teacher
	Department  string
	Categories  []string
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

func (c *BaseCourse) GetCode() string {
	return c.Code
}

func (c *BaseCourse) GetName() string {
	return c.Name
}

func (c *BaseCourse) GetCredit() float64 {
	return c.Credit
}

type Teacher struct {
	ID         int64
	Name       string
	Code       string
	Department string
	Title      string
}
