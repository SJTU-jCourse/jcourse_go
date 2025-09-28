package command

import (
	"context"

	"jcourse_go/internal/domain/reaction"
	"jcourse_go/internal/domain/review"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/pkg/apperror"
)

type ReactionService interface {
	CreateReaction(ctx context.Context, req shared.RequestCtx, cmd reaction.CreateReactionCommand) error
	DeleteReaction(ctx context.Context, req shared.RequestCtx, cmd reaction.DeleteReactionCommand) error
}

type reactionService struct {
	reactionRepo reaction.ReactionRepository
	reviewRepo   review.ReviewRepository
}

func (r *reactionService) CreateReaction(ctx context.Context, req shared.RequestCtx, cmd reaction.CreateReactionCommand) error {
	review, err := r.reviewRepo.Get(ctx, cmd.ReviewID)
	if err != nil {
		return err
	}
	if review == nil {
		return apperror.ErrNotFound
	}
	rct := reaction.NewReaction(cmd, req.User.UserID)
	return r.reactionRepo.Save(ctx, &rct)
}

func (r *reactionService) DeleteReaction(ctx context.Context, req shared.RequestCtx, cmd reaction.DeleteReactionCommand) error {
	rct, err := r.reactionRepo.Get(ctx, cmd.ReactionID)
	if err != nil || rct == nil {
		return err
	}
	if rct.UserID != req.User.UserID {
		return apperror.ErrNoPermission
	}
	return r.reactionRepo.Delete(ctx, cmd.ReactionID)
}

func NewReactionService(reactionRepo reaction.ReactionRepository, reviewRepo review.ReviewRepository) ReactionService {
	return &reactionService{
		reactionRepo: reactionRepo,
		reviewRepo:   reviewRepo,
	}
}
