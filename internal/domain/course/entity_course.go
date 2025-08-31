package course

import (
	"time"
)

// BaseCourse 基础课程，Code 作为 id
type BaseCourse struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Credit float64 `json:"credit"`
}

type Course struct {
	ID     string  `json:"id"`
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Credit float64 `json:"credit"`

	MainTeacherID string   `json:"main_teacher_id"`
	MainTeacher   *Teacher `json:"main_teacher"`

	OfferedCourses []OfferedCourse `json:"offered_courses"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OfferedCourse struct {
	ID           int64     `json:"id"`
	Semester     string    `json:"semester"`
	Department   string    `json:"department"`
	Language     string    `json:"language"`
	Grade        []string  `json:"grade"`
	Categories   []string  `json:"categories"`
	TeacherGroup []Teacher `json:"teacher_group"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Teacher struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`

	Title      string `json:"title"`
	Department string `json:"department"`
	Pinyin     Pinyin `json:"pinyin"`

	MainCourses []Course `json:"main_courses"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TrainingPlan struct {
	ID         int64                `json:"id"`
	Code       string               `json:"code"`
	MajorName  string               `json:"name"`
	EntryYear  string               `json:"entry_year"`
	Department string               `json:"department"`
	Degree     string               `json:"degree"`
	MajorClass string               `json:"major_class"`
	MinCredits float64              `json:"min_credits"`
	TotalYear  int64                `json:"total_year"`
	Courses    []TrainingPlanCourse `json:"courses"`
}

type TrainingPlanCourse struct {
	BaseCourse      BaseCourse `json:"base_course"`
	SuggestSemester string     `json:"suggest_semester"`
	Category        string     `json:"category"`
}
