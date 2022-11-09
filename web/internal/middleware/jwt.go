package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"simplegame.com/simplegame/common/jwtx"
	"simplegame.com/simplegame/web/internal/errors"
)

func JWTAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			c.Error(errors.Error(errors.Unauthorized, "unauthorized"))
			c.Abort()
			return
		}
		_, err := jwtx.ParseToken(auth)
		if err != nil && !strings.Contains(err.Error(), "expired") {
			c.Error(errors.Error(errors.Unauthorized, "unauthorized"))
			c.Abort()
			return
		}
		// TODO
		// check auth

		// check refresh_token
		refresh := c.Request.Header.Get("RefreshToken")
		refreshClaims, err := jwtx.ParseToken(refresh)
		if err != nil {
			c.Error(errors.Error(errors.Unauthorized, "unauthorized"))
			c.Abort()
			return
		}
		// generate new token
		newToken, err := jwtx.GenerateToken(refreshClaims.Info, time.Hour*2)
		if err != nil {
			c.Error(errors.Error(errors.Unauthorized, "unauthorized"))
			c.Abort()
			return
		}
		c.Header("newToken", newToken)
		c.Request.Header.Set("Authorization", newToken)
		c.Next()
	}
}
