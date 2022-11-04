package http_proxy_middleware

import (
	"gateway/common/response"
	"gateway/dao"
	"gateway/tools"
	"github.com/gin-gonic/gin"
)

func HTTPFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.FailMsg("service not found")
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		//统计项 1 全站 2 服务 3 租户
		totalCounter, err := tools.FlowCounterHandler.GetCounter(tools.FlowTotal)
		if err != nil {
			response.FailMsg(err.Error())

			c.Abort()
			return
		}
		totalCounter.Increase()

		//dayCount, _ := totalCounter.GetDayData(time.Now())
		//fmt.Printf("totalCounter qps:%v,dayCount:%v", totalCounter.QPS, dayCount)
		serviceCounter, err := tools.FlowCounterHandler.GetCounter(tools.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			response.FailMsg(err.Error())
			c.Abort()
			return
		}
		serviceCounter.Increase()

		//dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		//fmt.Printf("serviceCounter qps:%v,dayCount:%v", serviceCounter.QPS, dayServiceCount)
		c.Next()
	}
}
