package middleware

import (
	"fmt"
	"gateway/common/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"runtime/debug"
)

// RecoveryMiddleware RecoveryMiddleware捕获所有panic，并且返回错误信息
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//先做一下日志记录
				fmt.Println(string(debug.Stack()))
				logrus.WithField("_com_panic", map[string]interface{}{
					"error": fmt.Sprint(err),
					"stack": string(debug.Stack()),
				}).Info()
				response.FailMsg(fmt.Sprint(err))
				return
			}
		}()
		c.Next()
	}
}
