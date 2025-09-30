package middleware

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/config"
	"jcourse_go/internal/domain/auth"
)

const (
	CtxKeyUserSession string = "user_session"
)

const (
	CookieKeySession = "session_id"
)

func Session(redisConf config.RedisConfig, conf config.MiddlewareConfig) gin.HandlerFunc {
	store, err := sessions.NewRedisStore(10, "tcp", redisConf.Addr, redisConf.Password, []byte(conf.SessionSecret))
	if err != nil {
		panic(err)
	}
	return sessions.Sessions(CookieKeySession, store)

}

func GetCurrentUser(c *gin.Context) *auth.Session {
	s, exists := c.Get(CtxKeyUserSession)
	if !exists {
		return nil
	}
	return s.(*auth.Session)
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(CookieKeySession)
		if user == nil {
			c.JSON(http.StatusUnauthorized, nil)
			c.Abort()
		}
		c.Set(CtxKeyUserSession, user)
		c.Next()
	}
}
