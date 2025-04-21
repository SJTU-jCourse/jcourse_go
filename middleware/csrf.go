package middleware

/*
func CSRF() gin.HandlerFunc {
	key := util.GetCSRFSecret()
	csrfMd := csrf.Protect([]byte(key),
		csrf.Secure(!util.IsDebug()),
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
		token := csrf.Token(c.Request)
		c.SetCookie(constant.CSRFCookieKey, token, int((time.Hour * 24).Seconds()), "/", "", !util.IsDebug(), true)
	}
}
*/
