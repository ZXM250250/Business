package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"io/ioutil"
	"time"
)

// 请求进入日志
func RequestInLog(c *gin.Context) {
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	c.Set("startExecTime", time.Now())
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Write body back
	logrus.WithField("request_in", map[string]interface{}{
		"uri":    c.Request.RequestURI,
		"method": c.Request.Method,
		"args":   c.Request.PostForm,
		"body":   string(bodyBytes),
		"from":   c.ClientIP(),
	}).Info()
}

// 请求输出日志
func RequestOutLog(c *gin.Context) {
	// after request
	endExecTime := time.Now()
	response, _ := c.Get("response")
	st, _ := c.Get("startExecTime")

	startExecTime, _ := st.(time.Time)
	logrus.WithField("_com_request_out", map[string]interface{}{
		"uri":       c.Request.RequestURI,
		"method":    c.Request.Method,
		"args":      c.Request.PostForm,
		"from":      c.ClientIP(),
		"response":  response,
		"proc_time": endExecTime.Sub(startExecTime).Seconds(),
	}).Info()
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		//todo 优化点4
		RequestInLog(c)
		defer RequestOutLog(c)
		c.Next()
	}
}
