package middleware

import (
	"encoding/gob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/config"
	"jcourse_go/internal/constant"
	"jcourse_go/internal/domain/user"
	"jcourse_go/pkg/util"
)

func InitSession(r *gin.Engine) {
	secret := util.GetSessionSecret()
	store, err := sessions.NewRedisStore(10, "tcp", config.GetAppConfig().Redis.Addr, config.GetAppConfig().Redis.Password, []byte(secret))
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions(constant.CookieSessionKey, store))
	gob.Register(&user.UserDetail{})
}
