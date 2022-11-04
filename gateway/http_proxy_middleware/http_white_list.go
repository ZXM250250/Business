package http_proxy_middleware

import (
	"fmt"
	"gateway/common/response"
	"gateway/dao"
	"gateway/tools"
	"github.com/gin-gonic/gin"
	"strings"
)

//匹配接入方式 基于请求信息
func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.FailMsg("service not found")
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		iplist := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			iplist = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(iplist) > 0 {
			if !tools.InStringSlice(iplist, c.ClientIP()) {
				response.FailMsg(fmt.Sprintf("%s not in white ip list", c.ClientIP()))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
