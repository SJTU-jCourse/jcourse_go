package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/command"
	"jcourse_go/internal/domain/reaction"
	"jcourse_go/internal/domain/shared"
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
	var req reaction.CreateReactionCommand
	if ctx.ShouldBind(&req) != nil {
		WriteBadArgumentResponse(ctx)
		return
	}

	reqCtx := shared.NewRequestCtx(0, shared.UserRoleNormal)
	err := c.reactionService.CreateReaction(ctx, reqCtx, req)
	if err != nil {
		WriteErrorResponse(ctx, err)
		return
	}
	WriteDataResponse(ctx, nil)
}

func (c *ReviewReactionController) DeleteReaction(ctx *gin.Context) {
	var req reaction.DeleteReactionCommand
	if ctx.ShouldBind(&req) != nil {
		WriteBadArgumentResponse(ctx)
		return
	}

	reqCtx := shared.NewRequestCtx(0, shared.UserRoleNormal)
	err := c.reactionService.DeleteReaction(ctx, reqCtx, req)
	if err != nil {
		WriteErrorResponse(ctx, err)
		return
	}
	WriteDataResponse(ctx, nil)
}
