package ginrsp

import "github.com/gin-gonic/gin"

func Error(code int, msg string) gin.H {
	return Custom(code, msg, gin.H{})
}

func Succ(msg string, data interface{}) gin.H {
	return Custom(200, msg, data)
}

func Custom(code int, msg string, data interface{}) gin.H {
	return gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	}
}
