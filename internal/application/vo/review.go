package vo

import (
	"jcourse_go/internal/domain/course"
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

func NewReviewVO(r *course.Review) ReviewVO {
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
		c := NewCourseInReviewVO(r.Course)
		reviewVO.Course = &c
	}
	return reviewVO
}
