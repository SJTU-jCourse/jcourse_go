package vo

import (
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/infrastructure/entity"
)

type ReviewVO struct {
	ID        int64             `json:"id"`
	Course    *CourseInReviewVO `json:"course,omitempty"`
	Comment   string            `json:"comment"`
	Rating    int64             `json:"rating"`
	Score     string            `json:"score"`
	Semester  string            `json:"semester"`
	CreatedAt int64             `json:"created_at"`
	UpdatedAt int64             `json:"updated_at"`
}

func NewReviewVOFromDomain(r *course.Review) ReviewVO {
	reviewVO := ReviewVO{
		ID:        r.ID.Int64(),
		Course:    nil,
		Comment:   r.Comment,
		Rating:    r.Rating,
		Score:     r.Score,
		Semester:  r.Semester,
		CreatedAt: r.CreatedAt.Unix(),
		UpdatedAt: r.UpdatedAt.Unix(),
	}
	if r.Course != nil {
		c := NewCourseInReviewVOFromDomain(r.Course)
		reviewVO.Course = &c
	}
	return reviewVO
}

func NewReviewVOFromEntity(e *entity.Review) ReviewVO {
	courseVO := NewCourseInReviewVOFromEntity(e.Course)
	return ReviewVO{
		ID:        e.ID,
		Course:    &courseVO,
		Comment:   e.Comment,
		Rating:    e.Rating,
		Score:     e.Score,
		Semester:  e.Semester,
		CreatedAt: e.CreatedAt.Unix(),
		UpdatedAt: e.UpdatedAt.Unix(),
	}
}
