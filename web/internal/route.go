package internal

import (
	"net/http"
	"simplegame.com/simplegame/web/internal/ginrsp"

	"github.com/gin-gonic/gin"
)

func initRoute(engine *gin.Engine) {
	registerTestRoute(engine)
}

func registerTestRoute(engine *gin.Engine) {
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, ginrsp.Succ("pong", gin.H{}))
	})
}
