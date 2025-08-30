package middleware

import (
	"encoding/gob"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/constant"
	"jcourse_go/internal/dal"
	"jcourse_go/internal/model/model"
	"jcourse_go/pkg/util"
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
