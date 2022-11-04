package http_proxy_middleware

import (
	"gateway/common/response"
	"gateway/dao"
	"gateway/tools"
	"github.com/gin-gonic/gin"
	"strings"
)

//匹配接入方式 基于请求信息
func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.FailMsg("service not found")
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		whileIpList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			whileIpList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		blackIpList := []string{}
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(whileIpList) == 0 && len(blackIpList) > 0 {
			if tools.InStringSlice(blackIpList, c.ClientIP()) {
				response.FailMsg("%s in black ip list")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
