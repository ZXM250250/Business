package middleware

import (
	"gateway/common/response"
	"gateway/tools"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"net/http"
)

func JWTAuth(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusForbidden, response.FailMsg(response.FailureNullToken))
		c.Abort()
		return
	}
	logrus.Info("token=", token)
	user, err := tools.ParserToken(token)

	if err != nil {
		c.JSON(http.StatusForbidden, response.FailMsg(response.FailureParserToken))
		c.Abort()
		return
	}
	c.Set("user", user)

}
