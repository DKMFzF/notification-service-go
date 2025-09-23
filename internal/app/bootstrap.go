package app

import (
	"fmt"
	config "notification/internal/config"
	kafka "notification/internal/kafka"
	logger "notification/internal/logger"
	"notification/internal/middleware"
	api "notification/internal/routes"
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
	const (
		topic = "my-topic"
	)

	p, err := kafka.NewProducer([]string{app.Config.KafkaBroker})
	if err != nil {
		app.Logger.Fatalf("%s", "Error in creating producer", err)
	}

	for i := range 10 {
		msg := fmt.Sprintf("kafka message %d", i)
		if err := p.Produce(msg, topic); err != nil {
			app.Logger.Infof("Added msg '%s' in kafka for topic %s", msg, topic)
		}
	}

	app.Logger.Infof("Server bootstrap done")
	app.Logger.Infof("Server Port: %s", app.Config.Port)

	return *app
}
