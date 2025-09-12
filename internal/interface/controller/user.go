package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/query"
)

type UserController struct {
	userQuery query.UserQueryService
}

func NewUserController(
	userQuery query.UserQueryService,

) *UserController {
	return &UserController{
		userQuery: userQuery,
	}
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {

}

func (c *UserController) UpdateUserInfo(ctx *gin.Context) {

}
