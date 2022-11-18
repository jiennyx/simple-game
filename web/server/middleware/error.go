package middleware

import (
	"net/http"

	"simplegame.com/simplegame/web/server/errorx"
	"simplegame.com/simplegame/web/server/ginrsp"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		// check one
		if err, ok := c.Errors[0].Err.(*errorx.GinError); ok {
			c.JSON(http.StatusOK, ginrsp.Error(err.Code, err.Msg))
			return
		}
		c.JSON(http.StatusInternalServerError,
			ginrsp.Error(errorx.ServerError, "server error"))
	}
}
