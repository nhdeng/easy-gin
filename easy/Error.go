package easy

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	HTTP_STATUS = "Easy_STATUS"
)

func ErrorHandle() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": e,
					"status":  10001,
					"data":    nil,
				})
			}
		}()
		context.Next()
	}
}

func Throw(err string, code int, context *gin.Context) {
	context.Set(HTTP_STATUS, code)
	panic(err)
}

func Error(err error, msg ...string) {
	if err == nil {
		return
	} else {
		errMsg := err.Error()
		if len(msg) > 0 {
			errMsg = msg[0]
		}
		panic(errMsg)
	}

}
