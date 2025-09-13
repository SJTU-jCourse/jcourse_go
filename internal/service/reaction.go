package service

import (
	"context"
	"errors"

	"jcourse_go/internal/domain/user"
	"jcourse_go/internal/infrastructure/entity"
	"jcourse_go/internal/infrastructure/repository"
	"jcourse_go/internal/interface/dto"
)

func CreateReviewReaction(ctx context.Context, request olddto.CreateReviewReactionRequest, user *user.UserDetail) (int64, error) {

	if user == nil {
		return 0, errors.New("user not login")
	}

	r := repository.Q.ReviewPO
	_, err := r.WithContext(ctx).Where(r.ID.Eq(request.ReviewID)).Take()
	if err != nil {
		return 0, err
	}

	rq := repository.Q.ReviewReactionPO
	reactionModel := entity.ReviewReaction{
		ReviewID: request.ReviewID,
		UserID:   user.ID,
		Reaction: request.Reaction,
	}

	err = rq.WithContext(ctx).Create(&reactionModel)
	if err != nil {
		return 0, err
	}
	return reactionModel.ReviewID, nil
}

func DeleteReviewReaction(ctx context.Context, user *user.UserDetail, reactionID int64) error {
	r := repository.Q.ReviewReactionPO
	reaction, err := r.WithContext(ctx).Where(r.ID.Eq(reactionID)).Take()
	if err != nil {
		return err
	}
	if user != nil && reaction.UserID != user.ID {
		return errors.New("user not match")
	}
	_, err = r.WithContext(ctx).Delete(reaction)
	return err
}
