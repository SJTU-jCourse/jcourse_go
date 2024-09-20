package model

import "time"

type ReviewFilter struct {
	Page        int64
	PageSize    int64
	CourseID    int64
	Semester    string
	UserID      int64
	ReviewID    int64
	SearchQuery string
}

type Review struct {
	ID          int64         `json:"id"`
	Course      CourseMinimal `json:"course"`
	User        UserMinimal   `json:"user"`
	Comment     string        `json:"comment"`
	Rating      int64         `json:"rating"`
	Semester    string        `json:"semester"`
	IsAnonymous bool          `json:"is_anonymous"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
}
