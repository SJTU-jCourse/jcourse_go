package course

import "context"

type ReviewService interface {
	WriteReview(ctx context.Context) error
	UpdateReview(ctx context.Context) error
	DeleteReview(ctx context.Context) error
}
