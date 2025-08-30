package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/constant"
	dto2 "jcourse_go/internal/model/dto"
	"jcourse_go/internal/model/model"
	auth2 "jcourse_go/internal/service/auth"
)

type AuthController struct {
	authService *auth2.AuthService
}

func NewAuthController(authService *auth2.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (a *AuthController) LoginHandler(c *gin.Context) {
	var request dto2.LoginRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	user, err := a.authService.Login(c, request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "登录失败，请重试。"})
		return
	}
	err = storeAuthSession(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "登录失败，请重试。"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (a *AuthController) LogoutHandler(c *gin.Context) {
	clearAuthSession(c)
	c.JSON(http.StatusOK, dto2.BaseResponse{Message: "已登出"})
}

func (a *AuthController) ResetPasswordHandler(c *gin.Context) {
	var request dto2.ResetPasswordRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	err = a.authService.ResetPassword(c, request.Email, request.Password, request.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "重置密码失败，请重试。"})
		return
	}
	clearAuthSession(c)
	c.JSON(http.StatusOK, dto2.BaseResponse{Message: "重置密码成功"})
}

func (a *AuthController) RegisterHandler(c *gin.Context) {
	var request dto2.RegisterUserRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}
	user, err := a.authService.Register(c, request.Email, request.Password, request.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "注册失败，请重试。"})
		return
	}
	err = storeAuthSession(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "注册失败，请重试。"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (a *AuthController) SendVerifyCodeHandler(c *gin.Context) {
	var request dto2.SendEmailCodeRequest
	err := c.ShouldBind(&request)
	if err != nil || !auth2.ValidateEmail(request.Email) {
		c.JSON(http.StatusBadRequest, dto2.BaseResponse{Message: "参数错误"})
		return
	}

	err = a.authService.SendRegisterCodeEmail(c, request.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto2.BaseResponse{Message: "验证码发送失败，请稍后重试。"})
		return
	}

	c.JSON(http.StatusOK, dto2.BaseResponse{Message: "邮件已发送！请查看你的邮箱收件箱（包括垃圾邮件）"})
}

func clearAuthSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Set(constant.SessionUserAuthKey, nil)
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	_ = session.Save()
}

func storeAuthSession(c *gin.Context, user *model.UserDetail) error {
	if user == nil {
		return errors.New("user is nil")
	}
	session := sessions.Default(c)
	session.Set(constant.SessionUserAuthKey, user)
	err := session.Save()
	return err
}
