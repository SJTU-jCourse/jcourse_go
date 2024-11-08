package middleware

import (
	"jcourse_go/util"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"jcourse_go/constant"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
)

func GetCurrentUser(c *gin.Context) *model.UserDetail {
	user, exists := c.Get(constant.CtxKeyUser)
	if !exists {
		return nil
	}
	return user.(*model.UserDetail)
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
		userDomain, ok := user.(model.UserDetail)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.BaseResponse{Message: "未授权的请求"})
			c.Abort()
		}
		if userDomain.Role != model.UserRoleAdmin {
			c.JSON(http.StatusForbidden, dto.BaseResponse{Message: "未授权的请求"})
			c.Abort()
		}
		c.Set(constant.CtxKeyUser, user)
		c.Next()
	}
}

func IsMine(c *gin.Context, userID int64) bool {
	if util.IsNoLoginMode() {
		return true
	}
	user := GetCurrentUser(c)
	if user == nil {
		return false
	}
	return user.ID == userID
}
func IsAdmin(c *gin.Context) bool {
	if util.IsNoLoginMode() {
		return true
	}
	user := GetCurrentUser(c)
	if user == nil {
		return false
	}
	return user.Role == model.UserRoleAdmin
}

func IsMineOrAdmin(c *gin.Context, userID int64) bool {
	return IsMine(c, userID) || IsAdmin(c)
}
