package reaction

import (
	"context"

	"jcourse_go/internal/domain/shared"
)

type ReactionService interface {
	CreateReaction(ctx context.Context, req shared.RequestCtx, cmd CreateReactionCommand) error
	DeleteReaction(ctx context.Context, req shared.RequestCtx, cmd DeleteReactionCommand) error
}
