package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/command"
	"jcourse_go/internal/application/query"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/domain/user"
)

type UserController struct {
	userQuery          query.UserQueryService
	userProfileService command.UserProfileService
}

func NewUserController(
	userQuery query.UserQueryService,
	userProfileService command.UserProfileService,
) *UserController {
	return &UserController{
		userQuery:          userQuery,
		userProfileService: userProfileService,
	}
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID := shared.IDType(0)

	userInfo, err := c.userQuery.GetUserInfo(ctx, userID)
	if err != nil {
		WriteErrorResponse(ctx, err)
		return
	}
	WriteDataResponse(ctx, userInfo)
}

func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	var req user.UpdateUserInfoCommand
	if err := ctx.ShouldBind(&req); err != nil {
		WriteBadArgumentResponse(ctx)
		return
	}

	reqCtx := shared.NewRequestCtx(0, shared.UserRoleNormal)
	if err := c.userProfileService.UpdateUserInfo(ctx, reqCtx, req); err != nil {
		WriteErrorResponse(ctx, err)
		return
	}
	WriteDataResponse(ctx, nil)
}
