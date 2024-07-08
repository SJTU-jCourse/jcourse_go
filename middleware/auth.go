package middleware

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"jcourse_go/constant"
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
)

func GetUser(c *gin.Context) *domain.User {
	user, exists := c.Get(constant.CtxKeyUser)
	if !exists {
		return nil
	}
	return user.(*domain.User)
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(constant.SessionUserAuthKey)
		if user == nil {
			c.JSON(http.StatusUnauthorized, dto.BaseResponse{Message: "未授权的请求"})
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
			c.JSON(http.StatusUnauthorized, dto.BaseResponse{Message: "未授权的请求"})
			c.Abort()
		}
		userDomain, ok := user.(domain.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.BaseResponse{Message: "未授权的请求"})
			c.Abort()
		}
		if userDomain.Role != domain.UserRoleAdmin {
			c.JSON(http.StatusForbidden, dto.BaseResponse{Message: "未授权的请求"})
			c.Abort()
		}
		c.Set(constant.CtxKeyUser, user)
		c.Next()
	}
}
