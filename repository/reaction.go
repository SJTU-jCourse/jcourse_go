package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jcourse_go/model/po"
)

type IReviewReactionQuery interface {
	CreateReaction(ctx context.Context, reviewID, userID int64, reaction string) (*po.ReviewReactionPO, error)
	GetReaction(ctx context.Context, reactionID int64) (*po.ReviewReactionPO, error)
	DeleteReaction(ctx context.Context, reactionID int64) error
	GetReviewReactions(ctx context.Context, reviewIDs []int64) (map[int64][]po.ReviewReactionPO, error)
}

type ReviewReactionQuery struct {
	DB *gorm.DB
}

func (r *ReviewReactionQuery) GetReaction(ctx context.Context, reactionID int64) (*po.ReviewReactionPO, error) {
	reaction := po.ReviewReactionPO{}
	err := r.DB.WithContext(ctx).Model(&po.ReviewReactionPO{}).Where("id = ?", reactionID).First(&reaction).Error
	if err != nil {
		return nil, err
	}
	return &reaction, nil
}

func (r *ReviewReactionQuery) CreateReaction(ctx context.Context, reviewID, userID int64, reaction string) (*po.ReviewReactionPO, error) {
	reactionModel := po.ReviewReactionPO{
		ReviewID: reviewID,
		UserID:   userID,
		Reaction: reaction,
	}
	err := r.DB.WithContext(ctx).Model(&po.ReviewReactionPO{}).Create(&reactionModel).Error
	if err != nil {
		return nil, err
	}
	return &reactionModel, nil
}

func (r *ReviewReactionQuery) DeleteReaction(ctx context.Context, reactionID int64) error {
	reaction := po.ReviewReactionPO{}
	err := r.DB.WithContext(ctx).Model(&po.ReviewReactionPO{}).
		Where("id = ?", reactionID).Clauses(clause.Returning{}).Delete(&reaction).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewReactionQuery) GetReviewReactions(ctx context.Context, reviewIDs []int64) (map[int64][]po.ReviewReactionPO, error) {
	result := make(map[int64][]po.ReviewReactionPO)
	reactions := make([]po.ReviewReactionPO, 0)
	err := r.DB.WithContext(ctx).Model(&po.ReviewReactionPO{}).Where("review_id in ?", reviewIDs).Find(&reactions).Error
	if err != nil {
		return result, err
	}
	for _, reaction := range reactions {
		l, ok := result[reaction.ReviewID]
		if !ok {
			l = make([]po.ReviewReactionPO, 0)
		}
		l = append(l, reaction)
		result[reaction.ReviewID] = l
	}
	return result, nil
}

func NewReviewReactionQuery(db *gorm.DB) IReviewReactionQuery {
	return &ReviewReactionQuery{DB: db}
}
