package vo

import (
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/rating"
)

type CourseMinimal struct {
	course.BaseCourse
	ID          int64          `json:"id"`
	MainTeacher TeacherSummary `json:"main_teacher"`
}

type CourseSummary struct {
	CourseMinimal
	Categories []string          `json:"categories"`
	Department string            `json:"department"`
	RatingInfo rating.RatingInfo `json:"rating_info"`
}

type CourseDetail struct {
	CourseSummary
	OfferedCourse  []OfferedCourse `json:"offered_courses"`
	RelatedCourses *RelatedCourse  `json:"related_courses"`
}

type RelatedCourse struct {
	CoursesUnderSameTeacher  []CourseSummary `json:"courses_under_same_teacher"`
	CoursesWithOtherTeachers []CourseSummary `json:"courses_with_other_teachers"`
}
