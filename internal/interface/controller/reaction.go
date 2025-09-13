package controller

import (
	"github.com/gin-gonic/gin"

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

func (c *ReviewReactionController) CreateReaction(ctx *gin.Context) {
}

func (c *ReviewReactionController) DeleteReaction(ctx *gin.Context) {
}
