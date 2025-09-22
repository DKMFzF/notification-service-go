package middleware

import (
	pkgLogger "notification/pkg/logger"
	"strconv"

	gin "github.com/gin-gonic/gin"
)

func LoggerHandler(log pkgLogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Infof(
			"%s", "Request "+
				c.Request.RequestURI+
				" Response Code "+strconv.Itoa(c.Writer.Status()),
		)
	}
}
