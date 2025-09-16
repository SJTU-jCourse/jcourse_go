package reaction

import "jcourse_go/internal/domain/shared"

type CreateReactionCommand struct {
	Reaction string        `json:"reaction,omitempty"`
	ReviewID shared.IDType `json:"review_id,omitempty"`
	UserID   shared.IDType `json:"user_id,omitempty"`
}

type DeleteReactionCommand struct {
	ReactionID shared.IDType `json:"reaction_id,omitempty"`
}
