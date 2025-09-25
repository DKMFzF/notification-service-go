package routes

import (
	config "notification/internal/config"
	handlers "notification/internal/handlers"

	gin "github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) *gin.RouterGroup {
	apiGroup := router.Group("/")
	apiGroup.GET("/health", handlers.HealthHandler)
	return apiGroup
}
