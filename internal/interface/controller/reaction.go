package controller

import (
	"jcourse_go/internal/application/command"
)

type ReviewReactionController struct {
	reactionService command.ReactionService
}

func NewReviewReactionController(
	reactionService command.ReactionService,
) *ReviewReactionController {
	return &ReviewReactionController{
		reactionService: reactionService,
	}
}
