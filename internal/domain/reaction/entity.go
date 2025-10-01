package reaction

import (
	"time"

	"jcourse_go/internal/domain/shared"
)

type Reaction string

const (
	ReactionLike    Reaction = "like"
	ReactionDislike Reaction = "dislike"
)

type UserReaction struct {
	ID shared.IDType `json:"id"`

	ReviewID shared.IDType `json:"review_id"`
	UserID   shared.IDType `json:"user_id"`
	Reaction Reaction      `json:"reaction"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewUserReaction(cmd CreateReactionCommand, userID shared.IDType) UserReaction {
	return UserReaction{
		ReviewID:  cmd.ReviewID,
		UserID:    userID,
		Reaction:  cmd.Reaction,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
