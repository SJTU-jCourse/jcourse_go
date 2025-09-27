package controller

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/command"
	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/interface/dto"
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
	var req auth.LoginCommand
	if err := c.ShouldBind(&req); err != nil {
		dto.WriteBadArgumentResponse(c)
		return
	}
	if err := a.authService.Login(c, req); err != nil {
		dto.WriteErrorResponse(c, err)
		return
	}
	dto.WriteDataResponse(c, nil)
}

func (a *AuthController) LogoutHandler(c *gin.Context) {
	reqCtx := shared.NewRequestCtx(0, shared.UserRoleNormal)
	if err := a.authService.Logout(c, reqCtx); err != nil {
		dto.WriteErrorResponse(c, err)
		return
	}
	dto.WriteDataResponse(c, nil)
}

func (a *AuthController) ResetPasswordHandler(c *gin.Context) {
	var req auth.ResetPasswordCommand
	if err := c.ShouldBind(&req); err != nil {
		dto.WriteBadArgumentResponse(c)
		return
	}

	reqCtx := shared.NewRequestCtx(0, shared.UserRoleNormal)
	if err := a.authService.ResetPassword(c, reqCtx, req); err != nil {
		dto.WriteErrorResponse(c, err)
		return
	}
	dto.WriteDataResponse(c, nil)
}

func (a *AuthController) RegisterHandler(c *gin.Context) {
	var req auth.RegisterUserCommand
	if err := c.ShouldBind(&req); err != nil {
		dto.WriteBadArgumentResponse(c)
		return
	}

	if err := a.authService.Register(c, req); err != nil {
		dto.WriteErrorResponse(c, err)
		return
	}
	dto.WriteDataResponse(c, nil)
}

func (a *AuthController) SendVerifyCodeHandler(c *gin.Context) {
	var req auth.SendVerificationCodeCommand
	if err := c.ShouldBind(&req); err != nil {
		dto.WriteBadArgumentResponse(c)
		return
	}

	if err := a.authService.SendVerificationCode(c, req); err != nil {
		dto.WriteErrorResponse(c, err)
		return
	}
	dto.WriteDataResponse(c, nil)
}
