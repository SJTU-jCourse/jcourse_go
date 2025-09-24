package reaction

import (
	"context"

	"jcourse_go/internal/domain/shared"
)

type ReactionRepository interface {
	Get(ctx context.Context, reactionID shared.IDType) (*Reaction, error)
	Save(ctx context.Context, reaction *Reaction) error
	Delete(ctx context.Context, reactionID shared.IDType) error
}
