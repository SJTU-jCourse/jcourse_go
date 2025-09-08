package controller

import "jcourse_go/internal/application"

type ReviewReactionController struct {
	reactionService application.ReactionService
}

func NewReviewReactionController(
	reactionService application.ReactionService,
) *ReviewReactionController {
	return &ReviewReactionController{
		reactionService: reactionService,
	}
}
