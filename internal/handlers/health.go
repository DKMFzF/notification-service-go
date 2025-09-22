package handlers

import (
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
