package domain

type CourseLike interface {
	GetCode() string
	GetName() string
	GetCredit() float32
}

type OfferedCourse struct {
	Course
	ID           int64
	MainTeacher  Teacher
	TeacherGroup []Teacher
	Semester     string
	Department   string
	Categories   []string
	Location     string
	Language     string
	Grade        string
}

type Course struct {
	BaseCourse
	ID          int64
	MainTeacher Teacher
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
	Credit float32
}

func (c *BaseCourse) GetCode() string {
	return c.Code
}

func (c *BaseCourse) GetName() string {
	return c.Name
}

func (c *BaseCourse) GetCredit() float32 {
	return c.Credit
}

type Teacher struct {
	ID         int64
	Name       string
	Code       string
	Department string
	Title      string
}
