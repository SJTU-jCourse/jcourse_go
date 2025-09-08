package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application"
)

type UserController struct {
	userQuery application.UserQueryService
}

func NewUserController(
	userQuery application.UserQueryService,

) *UserController {
	return &UserController{
		userQuery: userQuery,
	}
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {

}

func (c *UserController) UpdateUserInfo(ctx *gin.Context) {

}
