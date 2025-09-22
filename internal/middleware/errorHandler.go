package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// central error handler
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errs := c.Errors
		if len(errs) > 0 {
			err := errs.Last()
			c.JSON(http.StatusInternalServerError, APIError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})

			c.Abort()
		}
	}
}
