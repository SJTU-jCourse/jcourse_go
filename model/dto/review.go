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
	Rating      int64             `json:"rating"`
	Semester    string            `json:"semester"`
	IsAnonymous bool              `json:"is_anonymous"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at,omitempty"`
}

type UpdateReviewDTO struct {
	ID          int64  `json:"id"`
	CourseID    int64  `json:"course_id" binding:"required"`
	Rating      int64  `json:"rating" binding:"required"`
	Comment     string `json:"comment" binding:"required"`
	Semester    string `json:"semester" binding:"required"`
	IsAnonymous bool   `json:"is_anonymous"`
}

type CreateReviewResponse struct {
	ReviewID int64 `json:"review_id"`
}

type ReviewListRequest struct {
	Page        int64  `json:"page" form:"page"`
	PageSize    int64  `json:"page_size" form:"page_size"`
	SearchQuery string `json:"search_query" form:"search_query"`
}

type ReviewListResponse = BasePaginateResponse[ReviewDTO]

type ReviewDetailRequest struct {
	ReviewID int64 `uri:"reviewID" binding:"required"`
}

type UpdateReviewRequest struct {
	ReviewID int64 `uri:"reviewID" binding:"required"`
}

type DeleteReviewRequest = UpdateReviewRequest

type UpdateReviewResponse = CreateReviewResponse

type DeleteReviewResponse = CreateReviewResponse
