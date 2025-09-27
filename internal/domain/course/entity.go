package course

import (
	"time"

	"jcourse_go/internal/domain/shared"
)

// Curriculum 基础课程，Code 作为 id
type Curriculum struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Credit float64 `json:"credit"`
}

type Course struct {
	ID     shared.IDType `json:"id"`
	Code   string        `json:"code"`
	Name   string        `json:"name"`
	Credit float64       `json:"credit"`

	MainTeacherID shared.IDType `json:"main_teacher_id"`
	MainTeacher   *Teacher      `json:"main_teacher"`

	OfferedCourses []CourseOffering `json:"offered_courses"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CourseOffering struct {
	ID           shared.IDType `json:"id"`
	Semester     string        `json:"semester"`
	Department   string        `json:"department"`
	Language     string        `json:"language"`
	Categories   []string      `json:"categories"`
	TeacherGroup []Teacher     `json:"teacher_group"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Teacher struct {
	ID   shared.IDType `json:"id"`
	Name string        `json:"name"`
	Code string        `json:"code"`

	Title      string `json:"title"`
	Department string `json:"department"`
	Pinyin     Pinyin `json:"pinyin"`

	Email   string `json:"email"`
	Picture string `json:"picture"`
	Bio     string `json:"bio"`

	MainCourses []Course `json:"main_courses"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TrainingPlan struct {
	ID          shared.IDType            `json:"id"`
	Code        string                   `json:"code"`
	MajorName   string                   `json:"name"`
	EntryYear   string                   `json:"entry_year"`
	Department  string                   `json:"department"`
	Degree      string                   `json:"degree"`
	MajorClass  string                   `json:"major_class"`
	MinCredits  float64                  `json:"min_credits"`
	TotalYear   int64                    `json:"total_year"`
	Curriculums []TrainingPlanCurriculum `json:"curriculums"`
}

type TrainingPlanCurriculum struct {
	Curriculum      Curriculum `json:"curriculum"`
	SuggestSemester string     `json:"suggest_semester"`
	Category        string     `json:"category"`
}
