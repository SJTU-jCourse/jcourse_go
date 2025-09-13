package vo

import (
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/infrastructure/entity"
)

type TeacherInCourseVO struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
}

func NewTeacherInCourseVOFromDomain(t *course.Teacher) TeacherInCourseVO {
	return TeacherInCourseVO{
		ID:         t.ID.Int64(),
		Name:       t.Name,
		Department: t.Department,
	}
}

func NewTeacherInCourseVOFromEntity(t *entity.Teacher) TeacherInCourseVO {
	return TeacherInCourseVO{
		ID:         t.ID,
		Name:       t.Name,
		Department: t.Department,
	}
}

type TeacherListItemVO struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Department string   `json:"department"`
	Title      string   `json:"title"`
	Picture    string   `json:"picture"`
	RatingInfo RatingVO `json:"rating_info"`
}

func NewTeacherListItemVO(t *course.Teacher) TeacherListItemVO {
	return TeacherListItemVO{
		ID:         t.ID.Int64(),
		Name:       t.Name,
		Department: t.Department,
		Title:      t.Title,
		Picture:    t.Picture,
	}
}

func NewTeacherListItemVOFromEntity(t *entity.Teacher) TeacherListItemVO {
	return TeacherListItemVO{
		ID:         t.ID,
		Name:       t.Name,
		Department: t.Department,
		Title:      t.Title,
		Picture:    t.Picture,
	}
}

type TeacherDetailVO struct {
	TeacherListItemVO
	Email string `json:"email"`
	Bio   string `json:"bio"`
}

func NewTeacherDetailVO(t *course.Teacher) TeacherDetailVO {
	return TeacherDetailVO{
		TeacherListItemVO: NewTeacherListItemVO(t),
		Email:             t.Email,
		Bio:               t.Bio,
	}
}

func NewTeacherDetailVOFromEntity(t *entity.Teacher) TeacherDetailVO {
	return TeacherDetailVO{
		TeacherListItemVO: NewTeacherListItemVOFromEntity(t),
		Email:             t.Email,
		Bio:               t.Biography,
	}
}
