package controller

import (
	"errors"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/constant"
	"jcourse_go/internal/domain/auth/model"
	"jcourse_go/internal/domain/auth/service"
	"jcourse_go/model/dto"
	"jcourse_go/pkg/util"
)

const (
	ResponseMessageLoginErr    = "登录失败，请稍后重试。"
	ResponseMessageRegisterErr = "注册失败，请稍后重试。"
	ResponseMessageSendCodeErr = "发送验证码失败，请稍后重试。"
	ResponseMessageResetPwdErr = "重置密码失败，请稍后重试。"
	ResponseMessageLogout      = "已登出。"
	ResponseMessageSendCode    = "邮件已发送！请查看你的邮箱收件箱（包括垃圾邮件）"
	ResponseMessageResetPwd    = "重置密码成功"
)

func RegisterAuthController(r *gin.RouterGroup, controller *AuthController) {
	auth := r.Group("/auth")
	auth.POST("/login", controller.Login)
	auth.POST("/logout", controller.Logout)
	auth.POST("/register", controller.Register)
	auth.POST("/send-verification-code", controller.SendVerificationCode)
	auth.POST("/reset-password", controller.ResetPassword)
}

type AuthController struct {
	authService service.AuthService
}

func (c *AuthController) Login(ctx *gin.Context) {
	var request dto.LoginRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		util.WrongParamResponse(ctx)
		return
	}
	user, err := c.authService.Login(ctx, request.Email, request.Password)
	if err != nil {
		util.ErrorResponse(ctx, ResponseMessageLoginErr)
		return
	}
	err = c.storeSession(ctx, user)
	if err != nil {
		util.ErrorResponse(ctx, ResponseMessageLoginErr)
		return
	}
	util.SuccessResponse(ctx, user)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	c.clearSession(ctx)
	util.SuccessSimpleResponse(ctx, ResponseMessageLogout)
}

func (c *AuthController) Register(ctx *gin.Context) {
	var request dto.RegisterUserRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		util.WrongParamResponse(ctx)
		return
	}

	user, err := c.authService.Register(ctx, request.Email, request.Password, request.Code)
	if err != nil {
		util.ErrorResponse(ctx, ResponseMessageRegisterErr)
		return
	}

	err = c.storeSession(ctx, user)
	if err != nil {
		util.ErrorResponse(ctx, ResponseMessageRegisterErr)
		return
	}
	util.SuccessResponse(ctx, user)
}

func (c *AuthController) SendVerificationCode(ctx *gin.Context) {
	var request dto.SendEmailCodeRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		util.WrongParamResponse(ctx)
		return
	}

	err = c.authService.SendVerificationCode(ctx, request.Email)
	if err != nil {
		util.ErrorResponse(ctx, ResponseMessageSendCodeErr)
		return
	}

	util.SuccessSimpleResponse(ctx, ResponseMessageSendCode)
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var request dto.ResetPasswordRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		util.WrongParamResponse(ctx)
		return
	}

	err = c.authService.ResetPassword(ctx, request.Email, request.Password, request.Code)
	if err != nil {
		util.ErrorResponse(ctx, ResponseMessageResetPwdErr)
		return
	}
	c.clearSession(ctx)
	util.SuccessSimpleResponse(ctx, ResponseMessageResetPwd)
}

func (c *AuthController) clearSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Set(constant.SessionUserAuthKey, nil)
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	_ = session.Save()
}

func (c *AuthController) storeSession(ctx *gin.Context, user *model.UserDomain) error {
	if user == nil {
		return errors.New("user is nil")
	}
	session := sessions.Default(ctx)
	session.Set(constant.SessionUserAuthKey, user)
	err := session.Save()
	return err
}
