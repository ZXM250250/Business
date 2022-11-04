package http_proxy_middleware

import (
	"fmt"
	"gateway/common/response"
	"gateway/dao"
	"gateway/tools"
	"github.com/gin-gonic/gin"
)

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.App)
		if appInfo.Qps > 0 {
			clientLimiter, err := tools.FlowLimiterHandler.GetLimiter(
				tools.FlowAppPrefix+appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.Qps))
			if err != nil {
				response.FailMsg(err.Error())
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				response.FailMsg(fmt.Sprintf("%v flow limit %v", c.ClientIP(), appInfo.Qps))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
