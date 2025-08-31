package course

import (
	"time"
)

type Review struct {
	ID          int64     `json:"id"`
	CourseID    int64     `json:"course_id"`
	Course      *Course   `json:"course"`
	UserID      int64     `json:"user_id"`
	Comment     string    `json:"comment"`
	Rating      int64     `json:"rating"`
	Semester    string    `json:"semester"`
	IsAnonymous bool      `json:"is_anonymous"`
	Grade       string    `json:"grade"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type ReviewRevision struct {
	ID          int64  `json:"id"`
	ReviewID    int64  `json:"review_id"`
	Comment     string `json:"comment"`
	Rating      int64  `json:"rating"`
	Semester    string `json:"semester"`
	IsAnonymous bool   `json:"is_anonymous"`
	Grade       string `json:"grade"`

	UpdatedBy int64     `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
