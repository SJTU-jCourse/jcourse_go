package olddto

type CreateReviewReactionRequest struct {
	ReviewID int64  `json:"review_id" binding:"required"`
	Reaction string `json:"reaction" binding:"required"`
}

type CreateReviewReactionResponse struct {
	ReactionID int64 `json:"reaction_id"`
}

type DeleteReviewReactionRequest struct {
	ReactionID int64 `json:"reaction_id" uri:"reactionID" binding:"required"`
}
