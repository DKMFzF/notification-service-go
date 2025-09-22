package routes

import (
	config "notification/internal/config"
	handlers "notification/internal/handlers"
	services "notification/internal/services"
	pkgLogger "notification/pkg/logger"

	gin "github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, log pkgLogger.Logger, cfg *config.Config) *gin.RouterGroup {
	apiGroup := router.Group("/")

	{
		emailsGroup := apiGroup.Group("emails")
		emailHandler := handlers.EmailHandler{
			Service: services.NewEmailService(cfg),
		}

		{
			emailsGroup.POST("/send", emailHandler.SendEmail)
		}
	}

	apiGroup.GET("/health", handlers.HealthHandler)

	return apiGroup
}
