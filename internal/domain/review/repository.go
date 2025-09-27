package review

import (
	"context"

	"jcourse_go/internal/domain/shared"
)

type ReviewRepository interface {
	Get(ctx context.Context, id shared.IDType) (*Review, error)
	Create(ctx context.Context, review *Review) error
	Update(ctx context.Context, review *Review, revision *ReviewRevision) error
	Delete(ctx context.Context, reviewID shared.IDType) error
}
