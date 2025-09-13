package middleware

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/constant"
	"jcourse_go/internal/domain/user"
	"jcourse_go/internal/interface/dto"
)

func GetCurrentUser(c *gin.Context) *user.UserDetail {
	user, exists := c.Get(constant.CtxKeyUser)
	if !exists {
		return nil
	}
	return user.(*user.UserDetail)
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(constant.SessionUserAuthKey)
		if user == nil {
			c.JSON(http.StatusUnauthorized, olddto.BaseResponse{Message: "未授权的请求"})
			c.Abort()
		}
		c.Set(constant.CtxKeyUser, user)
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(constant.SessionUserAuthKey)
		if user == nil {
			c.JSON(http.StatusUnauthorized, olddto.BaseResponse{Message: "未授权的请求"})
			c.Abort()
		}
		userDomain, ok := user.(user.UserDetail)
		if !ok {
			c.JSON(http.StatusUnauthorized, olddto.BaseResponse{Message: "未授权的请求"})
			c.Abort()
		}
		if userDomain.Role != user.UserRoleAdmin {
			c.JSON(http.StatusForbidden, olddto.BaseResponse{Message: "未授权的请求"})
			c.Abort()
		}
		c.Set(constant.CtxKeyUser, user)
		c.Next()
	}
}
