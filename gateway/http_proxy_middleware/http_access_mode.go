package http_proxy_middleware

import (
	"gateway/common/response"
	"gateway/dao"
	"github.com/gin-gonic/gin"
)

//匹配接入方式 基于请求信息
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			response.FailMsg(err.Error())

			c.Abort()
			return
		}
		//fmt.Println("matched service",tools.Obj2Json(service))
		c.Set("service", service)
		c.Next()
	}
}
