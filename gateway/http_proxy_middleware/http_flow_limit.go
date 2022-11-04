package http_proxy_middleware

import (
	"gateway/common/response"
	"gateway/dao"
	"gateway/tools"
	"github.com/gin-gonic/gin"
)

func HTTPFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.FailMsg("service not found")
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		if serviceDetail.AccessControl.ServiceFlowLimit != 0 {
			serviceLimiter, err := tools.FlowLimiterHandler.GetLimiter(
				tools.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				response.FailMsg(err.Error())

				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				response.FailMsg(err.Error())
				c.Abort()
				return
			}
		}

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := tools.FlowLimiterHandler.GetLimiter(
				tools.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+c.ClientIP(),
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				response.FailMsg(err.Error())

				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				response.FailMsg(err.Error())
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
