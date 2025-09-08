package application

import (
	"context"

	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
)

type ReviewCommandService interface {
	WriteReview(ctx context.Context, reqCtx shared.RequestCtx, cmd course.WriteReviewCommand) error
	UpdateReview(ctx context.Context, reqCtx shared.RequestCtx, cmd course.UpdateReviewCommand) error
	DeleteReview(ctx context.Context, reqCtx shared.RequestCtx, cmd course.DeleteReviewCommand) error
}

type reviewCommandService struct {
	courseRepo course.CourseRepository
	reviewRepo course.ReviewRepository
}

func (r *reviewCommandService) WriteReview(ctx context.Context, reqCtx shared.RequestCtx, cmd course.WriteReviewCommand) error {
	// TODO implement me
	panic("implement me")
}

func (r *reviewCommandService) UpdateReview(ctx context.Context, reqCtx shared.RequestCtx, cmd course.UpdateReviewCommand) error {
	// TODO implement me
	panic("implement me")
}

func (r *reviewCommandService) DeleteReview(ctx context.Context, reqCtx shared.RequestCtx, cmd course.DeleteReviewCommand) error {
	// TODO implement me
	panic("implement me")
}

func NewReviewCommandService(
	courseRepo course.CourseRepository,
	reviewRepo course.ReviewRepository,
) ReviewCommandService {
	return &reviewCommandService{
		courseRepo: courseRepo,
		reviewRepo: reviewRepo,
	}
}
