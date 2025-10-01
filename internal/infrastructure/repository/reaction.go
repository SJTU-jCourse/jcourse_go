package repository

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/reaction"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type ReactionRepository struct {
	db *gorm.DB
}

func (r *ReactionRepository) Get(ctx context.Context, reactionID shared.IDType) (*reaction.UserReaction, error) {
	e := &entity.ReviewReaction{}
	if err := r.db.Model(&entity.ReviewReaction{}).
		Where("id = ?", reactionID).
		First(&e).Error; err != nil {
		return nil, err
	}
	d := newReactionDomainFromEntity(e)
	return &d, nil
}

func (r *ReactionRepository) Save(ctx context.Context, reaction *reaction.UserReaction) error {
	e := newReactionEntityFromDomain(reaction)
	if err := r.db.Model(&entity.ReviewReaction{}).
		Create(&e).Error; err != nil {
		return err
	}
	reaction.ID = shared.IDType(e.ID)
	return nil
}

func (r *ReactionRepository) Delete(ctx context.Context, reactionID shared.IDType) error {
	if err := r.db.Delete(&entity.ReviewReaction{}, "id = ?", reactionID).Error; err != nil {
		return err
	}
	return nil
}

func NewReactionRepository(db *gorm.DB) reaction.ReactionRepository {
	return &ReactionRepository{db: db}
}

func newReactionDomainFromEntity(r *entity.ReviewReaction) reaction.UserReaction {
	rr, _ := reaction.ReactionString(r.Reaction)
	return reaction.UserReaction{
		ID:        shared.IDType(r.ID),
		ReviewID:  shared.IDType(r.ReviewID),
		UserID:    shared.IDType(r.UserID),
		Reaction:  rr,
		CreatedAt: r.CreatedAt,
	}
}

func newReactionEntityFromDomain(r *reaction.UserReaction) entity.ReviewReaction {
	return entity.ReviewReaction{
		ID:        int64(r.ID),
		ReviewID:  int64(r.ReviewID),
		UserID:    int64(r.UserID),
		Reaction:  r.Reaction.String(),
		CreatedAt: r.CreatedAt,
	}
}
