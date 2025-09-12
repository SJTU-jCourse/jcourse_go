package command

import (
	"context"

	"jcourse_go/internal/domain/reaction"
	"jcourse_go/internal/domain/shared"
)

type ReactionService interface {
	CreateReaction(ctx context.Context, req shared.RequestCtx, cmd reaction.CreateReactionCommand) error
	DeleteReaction(ctx context.Context, req shared.RequestCtx, cmd reaction.DeleteReactionCommand) error
}
