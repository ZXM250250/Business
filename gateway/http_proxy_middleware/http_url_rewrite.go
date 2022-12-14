package http_proxy_middleware

import (
	"gateway/common/response"
	"gateway/dao"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

//匹配接入方式 基于请求信息
func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.FailMsg("service not found")
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		for _, item := range strings.Split(serviceDetail.HTTPRule.UrlRewrite, ",") {
			//fmt.Println("item rewrite",item)
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}
			regexp, err := regexp.Compile(items[0])
			if err != nil {
				//fmt.Println("regexp.Compile err",err)
				continue
			}
			//fmt.Println("before rewrite",c.Request.URL.Path)
			replacePath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
			c.Request.URL.Path = string(replacePath)
			//fmt.Println("after rewrite",c.Request.URL.Path)
		}
		c.Next()
	}
}
