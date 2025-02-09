package service

import (
	"context"
	"errors"

	"jcourse_go/dal"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/repository"
)

func GetReviewReactionMap(ctx context.Context, reviewIDs []int64) (map[int64][]po.ReviewReactionPO, error) {
	reactionQuery := repository.NewReviewReactionQuery(dal.GetDBClient())
	reactionMap, err := reactionQuery.GetReviewReactions(ctx, reviewIDs)
	if err != nil {
		return nil, err
	}
	return reactionMap, nil
}

func CreateReviewReaction(ctx context.Context, request dto.CreateReviewReactionRequest, user *model.UserDetail) (int64, error) {

	if user == nil {
		return 0, errors.New("user not login")
	}

	reviewQuery := repository.NewReviewQuery(dal.GetDBClient())
	result, err := reviewQuery.GetReview(ctx, repository.WithID(request.ReviewID))
	if err != nil || len(result) == 0 {
		return 0, errors.New("review not exist")
	}

	reactionQuery := repository.NewReviewReactionQuery(dal.GetDBClient())
	reaction, err := reactionQuery.CreateReaction(ctx, request.ReviewID, user.ID, request.Reaction)
	if err != nil {
		return 0, err
	}
	return reaction.ReviewID, nil
}

func DeleteReviewReaction(ctx context.Context, user *model.UserDetail, reactionID int64) error {
	query := repository.NewReviewReactionQuery(dal.GetDBClient())
	reaction, err := query.GetReaction(ctx, reactionID)
	if err != nil {
		return err
	}
	if user != nil && reaction.UserID != user.ID {
		return errors.New("user not match")
	}
	return query.DeleteReaction(ctx, reactionID)
}
