package reaction

import "jcourse_go/internal/domain/shared"

type CreateReactionCommand struct {
	Reaction string
	ReviewID shared.IDType
	UserID   shared.IDType
}

type DeleteReactionCommand struct {
	ReactionID shared.IDType
}
