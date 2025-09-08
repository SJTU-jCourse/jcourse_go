package course

import "jcourse_go/internal/domain/shared"

type WriteReviewCommand struct {
	CourseID shared.IDType `json:"course_id,omitempty"`
	Comment  string        `json:"comment,omitempty"`
	Rating   int64         `json:"rating,omitempty"`
	Semester string        `json:"semester,omitempty"`
	Score    string        `json:"score,omitempty"`
}

type UpdateReviewCommand struct {
	ReviewID shared.IDType `json:"review_id,omitempty"`
	Comment  string        `json:"comment,omitempty"`
	Rating   int64         `json:"rating,omitempty"`
	Semester string        `json:"semester,omitempty"`
	Score    string        `json:"score,omitempty"`
}

type DeleteReviewCommand struct {
	ReviewID shared.IDType `json:"review_id,omitempty"`
}
