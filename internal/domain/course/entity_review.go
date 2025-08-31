package course

import (
	"time"

	"jcourse_go/internal/domain/shared"
)

type Review struct {
	ID          shared.IDType `json:"id"`
	CourseID    shared.IDType `json:"course_id"`
	Course      *Course       `json:"course"`
	UserID      shared.IDType `json:"user_id"`
	Comment     string        `json:"comment"`
	Rating      int64         `json:"rating"`
	Semester    string        `json:"semester"`
	IsAnonymous bool          `json:"is_anonymous"`
	Score       string        `json:"score"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
}

type ReviewRevision struct {
	ID          shared.IDType `json:"id"`
	ReviewID    shared.IDType `json:"review_id"`
	Comment     string        `json:"comment"`
	Rating      int64         `json:"rating"`
	Semester    string        `json:"semester"`
	IsAnonymous bool          `json:"is_anonymous"`
	Grade       string        `json:"grade"`

	UpdatedBy int64     `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
