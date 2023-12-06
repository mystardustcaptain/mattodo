package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mystardustcaptain/mattodo/pkg/config"
	"github.com/mystardustcaptain/mattodo/pkg/database"
	"github.com/mystardustcaptain/mattodo/pkg/route"
)

func main() {
	// Read configuration
	port := os.Getenv("SERVICE_PORT")

	// Initialize database
	db := database.InitDB(os.Getenv("DB_TYPE"), os.Getenv("DB_PATH"))

	// Initialize router
	r := route.InitializeRoutes(db)

	// Create a new server
	server := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %s\n", err)
		}
	}()

	// Channel to listen for interrupt or termination signals
	// SIGINT is sent when user presses Ctrl+C
	// SIGTERM is sent when you run kill command
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-stopChan
	log.Println("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the server, waiting max 10 seconds for current operations to complete
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %s\n", err)
	}

	log.Println("Server gracefully stopped")
}
