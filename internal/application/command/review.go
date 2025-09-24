package command

import (
	"context"
	"time"

	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/pkg/apperror"
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
	c, err := r.courseRepo.Get(ctx, cmd.CourseID)
	if err != nil {
		return err
	}
	if c == nil {
		return apperror.ErrNotFound
	}
	now := time.Now()
	rv, err := course.NewReview(cmd, reqCtx.User.UserID, now)
	if err != nil {
		return err
	}
	if err = r.reviewRepo.Create(ctx, &rv); err != nil {
		return err
	}
	return nil
}

func (r *reviewCommandService) UpdateReview(ctx context.Context, reqCtx shared.RequestCtx, cmd course.UpdateReviewCommand) error {
	rv, err := r.reviewRepo.Get(ctx, cmd.ReviewID)
	if err != nil {
		return err
	}
	if rv == nil {
		return apperror.ErrNotFound
	}
	if rv.UserID != reqCtx.User.UserID {
		return apperror.ErrNoPermission
	}
	now := time.Now()
	revision := rv.MakeRevision(reqCtx.User.UserID, now)
	if err = rv.BeUpdated(cmd, now); err != nil {
		return err
	}
	return r.reviewRepo.Update(ctx, rv, &revision)
}

func (r *reviewCommandService) DeleteReview(ctx context.Context, reqCtx shared.RequestCtx, cmd course.DeleteReviewCommand) error {
	rv, err := r.reviewRepo.Get(ctx, cmd.ReviewID)
	if err != nil || rv == nil {
		return err
	}
	if rv.UserID != reqCtx.User.UserID {
		return apperror.ErrNoPermission
	}
	return r.reviewRepo.Delete(ctx, cmd.ReviewID)
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
