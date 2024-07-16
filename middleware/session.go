package middleware

import (
	"encoding/gob"
	"os"

	"jcourse_go/constant"
	"jcourse_go/dal"
	"jcourse_go/model/domain"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InitSession(r *gin.Engine) {
	secret := os.Getenv(constant.SessionSecret)
	store, err := sessions.NewRedisStore(10, "tcp", dal.GetRedisDSN(), "", []byte(secret))
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions(constant.CookieSessionKey, store))
	gob.Register(&domain.User{})
}

func InitSessionDbg(r *gin.Engine) {
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions(constant.CookieSessionKey, store))
	gob.Register(&domain.User{})
}
