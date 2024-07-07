package dto

import "time"

type UserInReviewDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type ReviewDTO struct {
	ID          int64             `json:"id"`
	Course      CourseListItemDTO `json:"course"`
	User        UserInReviewDTO   `json:"user"`
	Comment     string            `json:"comment"`
	Rate        int64             `json:"rate"`
	Semester    string            `json:"semester"`
	IsAnonymous bool              `json:"is_anonymous"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at,omitempty"`
}

type CourseReviewWriteDTO struct {
	CourseID    int64  `json:"course_id"`
	Rate        int64  `json:"rate"`
	Comment     string `json:"comment"`
	Semester    string `json:"semester"`
	IsAnonymous bool   `json:"is_anonymous"`
}

type ReviewListRequest struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"page_size" form:"page_size"`
}

type ReviewListResponse = BasePaginateResponse[ReviewDTO]
