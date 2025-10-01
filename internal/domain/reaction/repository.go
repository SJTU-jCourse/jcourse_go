package reaction

import (
	"context"

	"jcourse_go/internal/domain/shared"
)

type ReactionRepository interface {
	Get(ctx context.Context, reactionID shared.IDType) (*UserReaction, error)
	Save(ctx context.Context, reaction *UserReaction) error
	Delete(ctx context.Context, reactionID shared.IDType) error
}
