package app

import (
	"context"
	http "net/http"
	config "notification/internal/config"
	handlers "notification/internal/handlers"
	kafka "notification/internal/kafka"
	logger "notification/internal/logger"
	middleware "notification/internal/middleware"
	api "notification/internal/routes"
	services "notification/internal/services"
	logType "notification/pkg/logger"
	signal "os/signal"
	"syscall"
	"time"

	kafkLib "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	gin "github.com/gin-gonic/gin"
)

type App struct {
	Server      *http.Server
	Context     context.Context
	Cancel      context.CancelFunc
	Router      *gin.Engine
	Logger      logType.Logger
	Config      *config.Config
	Consumer    *kafka.Consumer
	HandlersMap map[string]map[string]kafka.MessageHandler
}

func Bootstrap() *App {
	app := &App{
		Router: gin.New(),
		Logger: logger.Init(),
		Config: config.Load(),
	}
	app.Logger.Infof("App: INIT")

	app.Router.Use(gin.Recovery(), middleware.LoggerHandler(app.Logger), middleware.ErrorHandler())
	api.SetupRoutes(app.Router, app.Config)
	app.Logger.Infof("Routes: INIT")

	consumer, err := kafka.NewConsumer([]string{app.Config.KafkaBroker}, app.Config.KafkaGroupId)
	if err != nil {
		app.Logger.Fatalf("ERROR INIT consumer: %v", err)
	}
	app.Consumer = consumer
	app.Logger.Infof("Consumer: INIT")

	// list topics for subscribe
	topics := make([]string, 0, len(app.Config.KafkaTopics.Topics))
	for topic := range app.Config.KafkaTopics.Topics {
		topics = append(topics, topic)
	}
	if err := app.Consumer.Subscribe(topics); err != nil {
		app.Logger.Fatalf("Consumer not subscribe for kafka events: %v", err)
	}
	app.Logger.Infof("Subscribe on topics: INIT")

	// handlersMap for kafka
	app.HandlersMap = make(map[string]map[string]kafka.MessageHandler)
	for topic, events := range app.Config.KafkaTopics.Topics {
		app.HandlersMap[topic] = make(map[string]kafka.MessageHandler)

		for eventName, ev := range events {
			factory, ok := services.Get(ev.Service)
			if !ok {
				app.Logger.Errorf("unknown service: %s", ev.Service)
				continue
			}

			svc := factory.NewService(app.Config)
			conv := factory.Converter

			t, e := topic, eventName
			app.HandlersMap[t][e] = func(msg *kafkLib.Message) {
				h := handlers.NewKafkaMessageHandler(app.Logger, svc)
				h.HandleMessage(msg, conv)
			}
		}
	}
	app.Logger.Infof("Handlers map: INIT")

	// get started kafka
	app.Context, app.Cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	app.Consumer.Listen(app.Context, app.HandlersMap, app.Logger, app.Config.KafkaListenUpdate)
	app.Logger.Infof("Kafka consumer: INIT")

	return app
}

func (app *App) Run() {
	app.Server = &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: app.Router,
	}

	go func() {
		app.Logger.Infof("HTTP server starting on port: %s", app.Config.Port)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	// await signal finish
	<-app.Context.Done()
	app.Logger.Infof("Shutdown signal received")
	app.gracefulShutdown()
}

func (app *App) gracefulShutdown() {
	app.Cancel() // stop listen consumer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Server.Shutdown(ctx); err != nil {
		app.Logger.Errorf("Err down http server: %v", err)
	}

	app.Consumer.Close()
	app.Logger.Infof("Server close")
}
