package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"

	"jcourse_go/internal/config"
)

const (
	CookieKeyCSRF = "X-CSRF-Token"
)

func CSRF(conf config.MiddlewareConfig) gin.HandlerFunc {
	csrfMd := csrf.Protect([]byte(conf.CSRFSecret),
		csrf.Secure(conf.SecureCookie),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		})),
	)
	return adapter.Wrap(csrfMd)
}

func CSRFToken(conf config.MiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.Token(c.Request)
		c.SetCookie(CookieKeyCSRF, token, int((time.Hour * 24).Seconds()), "/", "", conf.SecureCookie, true)
	}
}

func CORS(conf config.MiddlewareConfig) gin.HandlerFunc {
	corsConf := cors.Config{
		AllowOrigins:     conf.CORSOrigin,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(corsConf)
}
