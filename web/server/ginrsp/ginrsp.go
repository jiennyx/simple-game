package ginrsp

import (
	"github.com/gin-gonic/gin"
	"simplegame.com/simplegame/web/server/errorx"
)

func FromError(err error) gin.H {
	switch e := err.(type) {
	case nil:
		return Succ(errorx.GetCodeMsg(errorx.OK), gin.H{})
	case *errorx.GinError:
		return Error(e.Code, e.Msg)
	default:
		return Error(errorx.ServerError, errorx.GetCodeMsg(errorx.ServerError))
	}
}

func Error(code int, msg string) gin.H {
	return Custom(code, msg, gin.H{})
}

func Succ(msg string, data any) gin.H {
	return Custom(errorx.OK, msg, data)
}

func Custom(code int, msg string, data any) gin.H {
	return gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	}
}
