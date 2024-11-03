package middleware

import (
	"encoding/gob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"jcourse_go/constant"
	"jcourse_go/dal"
	"jcourse_go/model/model"
	"jcourse_go/util"
)

func InitSession(r *gin.Engine) {
	secret := util.GetSessionSecret()
	store, err := sessions.NewRedisStore(10, "tcp", dal.GetRedisDSN(), dal.GetRedisPassWord(), []byte(secret))
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions(constant.CookieSessionKey, store))
	gob.Register(&model.UserDetail{})
}
