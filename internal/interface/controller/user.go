package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/domain/user"
)

type UserController struct {
	userRepo    user.UserRepository
	userService user.UserService
}

func NewUserController(
	userRepo user.UserRepository,
	userService user.UserService,
) *UserController {
	return &UserController{
		userRepo:    userRepo,
		userService: userService,
	}
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {

}

func (c *UserController) UpdateUserInfo(ctx *gin.Context) {

}
