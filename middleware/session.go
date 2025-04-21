package middleware

import (
	"encoding/gob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"jcourse_go/config"
	"jcourse_go/constant"
	"jcourse_go/dal"
	"jcourse_go/model/model"
)

func InitSession(conf *config.Config, r *gin.Engine) {
	secret := conf.Auth.SessionSecret
	store, err := sessions.NewRedisStore(10, "tcp", dal.GetRedisDSN(conf.Redis.Host, conf.Redis.Port), conf.Redis.Password, []byte(secret))
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions(constant.CookieSessionKey, store))
	gob.Register(&model.UserDetail{})
}
