package http_proxy_middleware

import (
	"gateway/common/response"
	"gateway/dao"
	"gateway/tools"
	"github.com/gin-gonic/gin"
	"strings"
)

//jwt auth token
func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.FailMsg("service not found")
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		//fmt.Println("serviceDetail",serviceDetail)
		// decode jwt token
		// app_id 与  app_list 取得 appInfo
		// appInfo 放到 gin.context
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		//fmt.Println("token",token)
		appMatched := false
		if token != "" {
			claims, err := tools.JwtDecode(token)
			if err != nil {
				response.FailMsg(err.Error())
				c.Abort()
				return
			}
			//fmt.Println("claims.Issuer",claims.Issuer)
			appList := dao.AppManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					c.Set("app", appInfo)
					appMatched = true
					break
				}
			}
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && !appMatched {
			response.FailMsg("not match valid app")
			c.Abort()
			return
		}
		c.Next()
	}
}
