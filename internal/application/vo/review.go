package vo

import (
	"jcourse_go/internal/domain/review"
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

func NewReviewVOFromDomain(r *review.Review) ReviewVO {
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
	return reviewVO
}

func NewReviewVOFromEntity(e *entity.Review) ReviewVO {
	reviewVO := ReviewVO{
		ID:        e.ID,
		Comment:   e.Comment,
		Rating:    e.Rating,
		Score:     e.Score,
		Semester:  e.Semester,
		CreatedAt: e.CreatedAt.Unix(),
		UpdatedAt: e.UpdatedAt.Unix(),
	}

	if e.Course != nil {
		courseVO := NewCourseInReviewVOFromEntity(e.Course)
		reviewVO.Course = &courseVO
	}

	return reviewVO
}
