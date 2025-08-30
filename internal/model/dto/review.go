package dto

import (
	model2 "jcourse_go/internal/model/model"
)

type UpdateReviewDTO struct {
	ID          int64  `json:"id"`
	CourseID    int64  `json:"course_id" binding:"required"`
	Rating      int64  `json:"rating" binding:"required"`
	Comment     string `json:"comment" binding:"required"`
	Semester    string `json:"semester" binding:"required"`
	IsAnonymous bool   `json:"is_anonymous"`
	Grade       string `json:"grade"`
}

type CreateReviewResponse struct {
	ReviewID int64 `json:"review_id"`
}

type ReviewListRequest struct {
	model2.PaginationFilterForQuery
	UserID   int64 `json:"user_id" form:"user_id"`
	CourseID int64 `json:"course_id" form:"course_id"`
	Rating   int64 `json:"rating" form:"rating"`
}

type ReviewListResponse = BasePaginateResponse[model2.Review]

type ReviewDetailRequest struct {
	ReviewID int64 `uri:"reviewID" binding:"required"`
}

type UpdateReviewRequest struct {
	ReviewID int64 `uri:"reviewID" binding:"required"`
}

type DeleteReviewRequest = UpdateReviewRequest

type UpdateReviewResponse = CreateReviewResponse

type DeleteReviewResponse = CreateReviewResponse
