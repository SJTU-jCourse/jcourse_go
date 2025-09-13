package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/command"
)

type AuthController struct {
	authService command.AuthService
}

func NewAuthController(authService command.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (a *AuthController) LoginHandler(c *gin.Context) {

}

func (a *AuthController) LogoutHandler(c *gin.Context) {

}

func (a *AuthController) ResetPasswordHandler(c *gin.Context) {

}

func (a *AuthController) RegisterHandler(c *gin.Context) {

}

func (a *AuthController) SendVerifyCodeHandler(c *gin.Context) {

}
