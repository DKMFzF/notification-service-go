package main

import (
	"context"
	"log"
	http "net/http"
	start "notification/internal/app"
	signal "os/signal"
	"syscall"
	"time"
)

// soft kill server process
func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// trigger kill process
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	// new context for 5s timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func main() {
	app := start.Bootstrap()

	srv := &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: app.Router,
	}

	// create check rutine
	done := make(chan bool, 1)
	go gracefulShutdown(srv, done)

	app.Logger.Infof("Start server...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	<-done
	app.Logger.Infof("Server exiting")
}
