package http_proxy_middleware

import (
	"fmt"
	"gateway/common/response"
	"gateway/dao"
	"gateway/tools"
	"github.com/gin-gonic/gin"
)

func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.App)
		appCounter, err := tools.FlowCounterHandler.GetCounter(tools.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			response.FailMsg(err.Error())
			c.Abort()
			return
		}
		appCounter.Increase()
		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			response.FailMsg(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", appInfo.Qpd, appCounter.TotalCount))
			c.Abort()
			return
		}
		c.Next()
	}
}
