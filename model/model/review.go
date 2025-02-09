package model

import "time"

type ReviewFilterForQuery struct {
	PaginationFilterForQuery
	CourseID         int64
	Semester         string
	UserID           int64
	ReviewID         int64
	Rating           int64
	ExcludeAnonymous bool
}

type Review struct {
	ID          int64          `json:"id"`
	Course      CourseMinimal  `json:"course"`
	User        UserMinimal    `json:"user"`
	Comment     string         `json:"comment"`
	Rating      int64          `json:"rating"`
	Semester    string         `json:"semester"`
	IsAnonymous bool           `json:"is_anonymous"`
	Grade       string         `json:"grade"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	Reaction    ReviewReaction `json:"reaction"`
}

type ReviewReaction struct {
	TotalReactions map[string]int64 `json:"total_reactions"` // reaction -> count
	MyReactions    map[string]int64 `json:"my_reactions"`    // reaction -> id
}
