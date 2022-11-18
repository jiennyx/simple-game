package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"simplegame.com/simplegame/common/jwtx"
	"simplegame.com/simplegame/web/server/errorx"
)

var (
	whiteRoute = map[string]bool{
		"/register": true,
		"/auth":     true,
	}
)

func isWhiteRoute(route string) bool {
	return whiteRoute[route]
}

// TODO
func JWTAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isWhiteRoute(c.Request.URL.Path) {
			c.Next()
			return
		}
		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			c.Error(errorx.Error(errorx.Unauthorized, "unauthorized: len(auth) == 0"))
			c.Abort()
			return
		}
		authClaims, err := jwtx.ParseToken(jwtx.JWTTypeAuth, auth)
		if err == nil {
			c.Next()
			return
		}
		if err != nil && !strings.Contains(err.Error(), "expired") {
			// c.Error(errorx.Error(errorx.Unauthorized, "unauthorized: authClaims"))
			c.Error(errorx.Error(errorx.Unauthorized, err.Error()))
			c.Abort()
			return
		}
		// check auth
		if !jwtx.IsValidIssuer(authClaims.Issuer) {
			c.Error(errorx.Error(errorx.Unauthorized, "unauthorized: not valid issuer"))
			c.Abort()
			return
		}

		// check refresh_token
		refresh := c.Request.Header.Get("RefreshToken")
		refreshClaims, err := jwtx.ParseToken(jwtx.JWTTypeRefresh, refresh)
		if err != nil {
			// c.Error(errorx.Error(errorx.Unauthorized, "unauthorized: refreshClaims"))
			c.Error(errorx.Error(errorx.Unauthorized, err.Error()))
			c.Abort()
			return
		}
		// generate new token
		// jwt info test
		newToken, err := jwtx.GenerateToken(
			jwtx.JWTTypeAuth,
			refreshClaims.Info,
			time.Minute*10,
		)
		if err != nil {
			c.Error(errorx.Error(errorx.Unauthorized, "unauthorized: generate failed"))
			c.Abort()
			return
		}
		c.Header("newToken", newToken)
		c.Request.Header.Set("Authorization", newToken)
		c.Next()
	}
}
