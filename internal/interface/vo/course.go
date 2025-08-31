package vo

import (
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/rating"
)

type OfferedInfoVO struct {
	Categories []string `json:"categories"`
	Department string   `json:"department"`
	Language   string   `json:"language"`
	Grade      []string `json:"grade"`
}

type CourseListItemVO struct {
	ID                int64             `json:"id"`
	Code              string            `json:"code"`
	Name              string            `json:"name"`
	Credit            float64           `json:"credit"`
	MainTeacher       TeacherInCourseVO `json:"main_teacher"`
	LatestOfferedInfo OfferedInfoVO     `json:"latest_offered_info"`
	RatingInfo        rating.RatingInfo `json:"rating_info"`
}

type CourseInReviewVO struct {
	ID          int64             `json:"id,omitempty"`
	Code        string            `json:"code,omitempty"`
	Name        string            `json:"name,omitempty"`
	MainTeacher TeacherInCourseVO `json:"main_teacher"`
}

type CourseDetailVO struct {
	ID          int64              `json:"id"`
	Code        string             `json:"code"`
	Name        string             `json:"name"`
	Credit      float64            `json:"credit"`
	MainTeacher TeacherInCourseVO  `json:"main_teacher"`
	Offering    []CourseOfferingVO `json:"offering"`
}

type CourseOfferingVO struct {
	Semester     string              `json:"semester"`
	TeacherGroup []TeacherInCourseVO `json:"teacher_group"`

	Categories []string `json:"categories"`
	Department string   `json:"department"`
	Language   string   `json:"language"`
	Grade      []string `json:"grade"`
}

func NewCourseOfferingVO(co *course.CourseOffering) CourseOfferingVO {
	teacherGroup := make([]TeacherInCourseVO, 0)
	for _, t := range co.TeacherGroup {
		teacher := NewTeacherInCourseVO(&t)
		teacherGroup = append(teacherGroup, teacher)
	}
	return CourseOfferingVO{
		Semester:     co.Semester,
		TeacherGroup: teacherGroup,
		Categories:   co.Categories,
		Department:   co.Department,
		Language:     co.Language,
		Grade:        co.Grade,
	}
}

type CurriculumVO struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Credit float64 `json:"credit"`
}

func NewCurriculumVO(c *course.Curriculum) CurriculumVO {
	return CurriculumVO{
		Code:   c.Code,
		Name:   c.Name,
		Credit: c.Credit,
	}
}
