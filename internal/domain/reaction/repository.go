package reaction

import (
	"context"

	"jcourse_go/internal/domain/shared"
)

type ReactionRepository interface {
	FindByReviewID(ctx context.Context, reviewID shared.IDType) ([]Reaction, error)
	FindByUserID(ctx context.Context, userID shared.IDType) ([]Reaction, error)
	Save(ctx context.Context, reaction *Reaction) error
	Delete(ctx context.Context, reaction *Reaction) error
}
