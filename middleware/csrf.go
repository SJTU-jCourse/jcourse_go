package middleware

import (
	"net/http"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"jcourse_go/constant"
	"jcourse_go/model/dto"
	"jcourse_go/util"
)

func CSRF() gin.HandlerFunc {
	key := os.Getenv(constant.CSRFSecretKey)
	csrfMd := csrf.Protect([]byte(key),
		csrf.MaxAge(0),
		csrf.Secure(util.IsDebug()),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			bytes, _ := sonic.Marshal(dto.BaseResponse{Message: "Forbidden - CSRF token invalid"})
			_, _ = w.Write(bytes)
		})),
	)
	return adapter.Wrap(csrfMd)
}

func CSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(constant.CSRFHeader, csrf.Token(c.Request))
	}
}
