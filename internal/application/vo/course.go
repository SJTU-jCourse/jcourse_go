package vo

import (
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/infrastructure/entity"
)

type OfferingInfoVO struct {
	Categories []string `json:"categories"`
	Department string   `json:"department"`
	Language   string   `json:"language"`
}

func NewOfferingInfoVOFromDomain(co *course.CourseOffering) OfferingInfoVO {
	return OfferingInfoVO{
		Categories: co.Categories,
		Department: co.Department,
		Language:   co.Language,
	}
}

func NewOfferingInfoVOFromEntity(co *entity.CourseOffering) OfferingInfoVO {
	categories := make([]string, 0)
	for _, ct := range co.Categories {
		categories = append(categories, ct.Category)
	}
	return OfferingInfoVO{
		Categories: categories,
		Department: co.Department,
		Language:   co.Language,
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

func NewCourseListItemVOFromDomain(c *course.Course) CourseListItemVO {
	var offering OfferingInfoVO
	if len(c.OfferedCourses) > 0 {
		offering = NewOfferingInfoVOFromDomain(&c.OfferedCourses[len(c.OfferedCourses)-1])
	}
	return CourseListItemVO{
		ID:             c.ID.Int64(),
		Code:           c.Code,
		Name:           c.Name,
		Credit:         c.Credit,
		MainTeacher:    NewTeacherInCourseVOFromDomain(c.MainTeacher),
		LatestOffering: offering,
	}
}

func NewCourseListItemVOFromEntity(e *entity.Course) CourseListItemVO {
	var offering OfferingInfoVO
	if len(e.Offerings) > 0 {
		offering = NewOfferingInfoVOFromEntity(e.Offerings[0])
	}
	return CourseListItemVO{
		ID:             e.ID,
		Code:           e.Code,
		Name:           e.Name,
		Credit:         e.Credit,
		MainTeacher:    NewTeacherInCourseVOFromEntity(e.MainTeacher),
		LatestOffering: offering,
	}
}

type CourseInReviewVO struct {
	ID          int64             `json:"id,omitempty"`
	Code        string            `json:"code,omitempty"`
	Name        string            `json:"name,omitempty"`
	MainTeacher TeacherInCourseVO `json:"main_teacher"`
}

func NewCourseInReviewVOFromDomain(c *course.Course) CourseInReviewVO {
	return CourseInReviewVO{
		ID:          c.ID.Int64(),
		Code:        c.Code,
		Name:        c.Name,
		MainTeacher: NewTeacherInCourseVOFromDomain(c.MainTeacher),
	}
}

func NewCourseInReviewVOFromEntity(e *entity.Course) CourseInReviewVO {
	return CourseInReviewVO{
		ID:          e.ID,
		Code:        e.Code,
		Name:        e.Name,
		MainTeacher: NewTeacherInCourseVOFromEntity(e.MainTeacher),
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

func NewCourseDetailVOFromDomain(c *course.Course) CourseDetailVO {
	offering := make([]CourseOfferingVO, 0)
	for _, o := range c.OfferedCourses {
		offering = append(offering, NewCourseOfferingVO(&o))
	}
	return CourseDetailVO{
		ID:          c.ID.Int64(),
		Code:        c.Code,
		Name:        c.Name,
		Credit:      c.Credit,
		MainTeacher: NewTeacherInCourseVOFromDomain(c.MainTeacher),
		Offering:    offering,
	}
}

func NewCourseDetailVOFromEntity(e *entity.Course) CourseDetailVO {
	offeringVOs := make([]CourseOfferingVO, 0)
	for _, co := range e.Offerings {
		coVO := NewCourseOfferingVOFromEntity(co)
		offeringVOs = append(offeringVOs, coVO)
	}
	return CourseDetailVO{
		ID:          e.ID,
		Code:        e.Code,
		Name:        e.Name,
		Credit:      e.Credit,
		MainTeacher: NewTeacherInCourseVOFromEntity(e.MainTeacher),
		Offering:    offeringVOs,
	}
}

type CourseOfferingVO struct {
	Semester     string              `json:"semester"`
	TeacherGroup []TeacherInCourseVO `json:"teacher_group"`

	Categories []string `json:"categories"`
	Department string   `json:"department"`
	Language   string   `json:"language"`
}

func NewCourseOfferingVO(co *course.CourseOffering) CourseOfferingVO {
	teacherGroup := make([]TeacherInCourseVO, 0)
	for _, t := range co.TeacherGroup {
		teacher := NewTeacherInCourseVOFromDomain(&t)
		teacherGroup = append(teacherGroup, teacher)
	}
	return CourseOfferingVO{
		Semester:     co.Semester,
		TeacherGroup: teacherGroup,
		Categories:   co.Categories,
		Department:   co.Department,
		Language:     co.Language,
	}
}

func NewCourseOfferingVOFromEntity(co *entity.CourseOffering) CourseOfferingVO {
	teacherGroup := make([]TeacherInCourseVO, 0)
	for _, t := range co.TeacherGroup {
		teacher := NewTeacherInCourseVOFromEntity(t)
		teacherGroup = append(teacherGroup, teacher)
	}
	categories := make([]string, 0)
	for _, ct := range co.Categories {
		categories = append(categories, ct.Category)
	}
	return CourseOfferingVO{
		Semester:     co.Semester,
		TeacherGroup: teacherGroup,
		Categories:   categories,
		Department:   co.Department,
		Language:     co.Language,
	}
}
