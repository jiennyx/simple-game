package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simplegame.com/simplegame/web/logic"
	"simplegame.com/simplegame/web/model"
	"simplegame.com/simplegame/web/server/errorx"
	"simplegame.com/simplegame/web/server/ginrsp"
)

func Register(c *gin.Context) {
	var reqObj model.RegisterReq
	err := c.ShouldBind(&reqObj)
	if err != nil {
		c.JSON(
			http.StatusOK,
			ginrsp.Error(errorx.ParamError, errorx.GetCodeMsg(errorx.ParamError)),
		)
		return
	}

	c.JSON(http.StatusOK, ginrsp.Succ("succ", gin.H{}))
}

func GetAuth(c *gin.Context) {
	var reqObj model.GetAuthReq
	err := c.ShouldBind(&reqObj)
	if err != nil {
		c.JSON(
			http.StatusOK,
			ginrsp.Error(errorx.ParamError, errorx.GetCodeMsg(errorx.ParamError)),
		)
		return
	}

	rsp, err := logic.GetAuth(c, reqObj)
	if err != nil {
		c.JSON(http.StatusOK, ginrsp.FromError(err))
		return
	}

	c.JSON(http.StatusOK, ginrsp.Succ("succ", rsp))
}
