package vo

import (
	"jcourse_go/internal/domain/course"
)

type OfferingInfoVO struct {
	Categories []string `json:"categories"`
	Department string   `json:"department"`
	Language   string   `json:"language"`
	Grade      []string `json:"grade"`
}

func NewOfferingInfoVO(co *course.CourseOffering) OfferingInfoVO {
	return OfferingInfoVO{
		Categories: co.Categories,
		Department: co.Department,
		Language:   co.Language,
		Grade:      co.Grade,
	}
}

type CourseListItemVO struct {
	ID             int64             `json:"id"`
	Code           string            `json:"code"`
	Name           string            `json:"name"`
	Credit         float64           `json:"credit"`
	MainTeacher    TeacherInCourseVO `json:"main_teacher"`
	LatestOffering OfferingInfoVO    `json:"latest_offering"`
	RatingInfo     RatingVO          `json:"rating_info"`
}

func NewCourseListItemVO(c *course.Course) CourseListItemVO {
	var offering OfferingInfoVO
	if len(c.OfferedCourses) > 0 {
		offering = NewOfferingInfoVO(&c.OfferedCourses[len(c.OfferedCourses)-1])
	}
	return CourseListItemVO{
		ID:             c.ID.Int64(),
		Code:           c.Code,
		Name:           c.Name,
		Credit:         c.Credit,
		MainTeacher:    NewTeacherInCourseVO(c.MainTeacher),
		LatestOffering: offering,
	}
}

type CourseInReviewVO struct {
	ID          int64             `json:"id,omitempty"`
	Code        string            `json:"code,omitempty"`
	Name        string            `json:"name,omitempty"`
	MainTeacher TeacherInCourseVO `json:"main_teacher"`
}

func NewCourseInReviewVO(c *course.Course) CourseInReviewVO {
	return CourseInReviewVO{
		ID:          c.ID.Int64(),
		Code:        c.Code,
		Name:        c.Name,
		MainTeacher: NewTeacherInCourseVO(c.MainTeacher),
	}
}

type CourseDetailVO struct {
	ID          int64              `json:"id"`
	Code        string             `json:"code"`
	Name        string             `json:"name"`
	Credit      float64            `json:"credit"`
	MainTeacher TeacherInCourseVO  `json:"main_teacher"`
	Offering    []CourseOfferingVO `json:"offering"`
}

func NewCourseDetailVO(c *course.Course) CourseDetailVO {
	offering := make([]CourseOfferingVO, 0)
	for _, o := range c.OfferedCourses {
		offering = append(offering, NewCourseOfferingVO(&o))
	}
	return CourseDetailVO{
		ID:          c.ID.Int64(),
		Code:        c.Code,
		Name:        c.Name,
		Credit:      c.Credit,
		MainTeacher: NewTeacherInCourseVO(c.MainTeacher),
		Offering:    offering,
	}
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
