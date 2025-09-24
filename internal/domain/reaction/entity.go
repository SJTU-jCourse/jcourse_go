package reaction

import (
	"time"

	"jcourse_go/internal/domain/shared"
)

type Reaction struct {
	ID shared.IDType `json:"id"`

	ReviewID shared.IDType `json:"review_id"`
	UserID   shared.IDType `json:"user_id"`
	Reaction string        `json:"reaction"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewReaction(cmd CreateReactionCommand, userID shared.IDType) Reaction {
	return Reaction{
		ReviewID:  cmd.ReviewID,
		UserID:    userID,
		Reaction:  cmd.Reaction,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
