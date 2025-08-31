package course

import (
	"context"

	"jcourse_go/internal/domain/shared"
)

type ReviewCommandService interface {
	WriteReview(ctx context.Context, req shared.RequestCtx, cmd WriteReviewCommand) error
	UpdateReview(ctx context.Context, req shared.RequestCtx, cmd UpdateReviewCommand) error
	DeleteReview(ctx context.Context, req shared.RequestCtx, cmd DeleteReviewCommand) error
}
