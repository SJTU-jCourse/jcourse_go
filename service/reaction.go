package service

import (
	"context"
	"errors"

	"jcourse_go/internal/infra/query"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

func CreateReviewReaction(ctx context.Context, request dto.CreateReviewReactionRequest, user *model.UserDetail) (int64, error) {

	if user == nil {
		return 0, errors.New("user not login")
	}

	r := query.Q.ReviewPO
	_, err := r.WithContext(ctx).Where(r.ID.Eq(request.ReviewID)).Take()
	if err != nil {
		return 0, err
	}

	rq := query.Q.ReviewReactionPO
	reactionModel := po.ReviewReactionPO{
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

func DeleteReviewReaction(ctx context.Context, user *model.UserDetail, reactionID int64) error {
	r := query.Q.ReviewReactionPO
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
