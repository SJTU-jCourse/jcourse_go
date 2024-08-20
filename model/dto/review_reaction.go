package dto

type ReviewReactionEmbedDTO struct {
	Reaction string `json:"reaction"`
	Count    uint   `json:"count"`
	Yours    bool   `json:"yours"`
}

// type ReviewReactionEmbedList = []ReviewReactionEmbedDTO

type CreateReviewReactionRequest struct {
	ReviewID int64  `json:"reviewID" binding:"required"`
	Reaction string `json:"reaction" binding:"required"`
}

type DeleteReviewReactionRequest struct {
	ReactionID int64 `uri:"reactionID" binding:"required"`
}

type CreateReviewReactionResponse struct {
	ReactionID int64 `json:"reaction_id"`
}

type DeleteReviewReactionResponse = CreateReviewReactionResponse
