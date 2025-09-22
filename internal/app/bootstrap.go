package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	config "notification/internal/config"
	kafka "notification/internal/kafka"
	logger "notification/internal/logger"
	"notification/internal/middleware"
	api "notification/internal/routes"
	services "notification/internal/services"
	logType "notification/pkg/logger"

	"github.com/gin-gonic/gin"
)

// global type for Application
type App struct {
	Router *gin.Engine
	Logger logType.Logger
	Config *config.Config
	Err    error
}

func Bootstrap() App {
	app := &App{
		Router: gin.New(),
		Logger: logger.Init(),
		Config: config.Load(),
	}

	// rest
	app.Router.Use(gin.Recovery())
	app.Router.Use(middleware.LoggerHandler(app.Logger))
	app.Router.Use(middleware.ErrorHandler())
	api.SetupRoutes(app.Router, app.Logger, app.Config)

	// kafka
	topics := []string{"user_registered", "user_deleted"}
	kafkaConsumer := kafka.NewNotificationConsumer(
		app.Config.KafkaBroker,
		"notification_group",
		topics,
		services.NewEmailService(app.Config),
		kafka.NewKafkaProducer(app.Config.KafkaBroker),
		app.Logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	kafkaConsumer.Start(ctx)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		app.Logger.Infof("shutdown signal stopping consumers...")
		cancel()
	}()

	app.Logger.Infof("Server bootstrap done")
	app.Logger.Infof("Server Port: %s", app.Config.Port)

	return *app
}
