package initialize

import (
	"fmt"
	"mlvt/internal/infra/zap-logging/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Run initializes the application and starts the server.
// It encapsulates all initialization logic and handles graceful shutdown.
func Run() {
	// Initialize Logger
	if err := InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Logger initialization failed: %v\n", err)
		os.Exit(1)
	}

	// Initialize Database
	dbConn, err := InitDatabase()
	if err != nil {
		log.Errorf("Database initialization failed: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Warnf("Error closing database connection: %v", err)
		}
	}()

	// Initialize AWS Clients
	// s3Client, err := InitAWS()
	// if err != nil {
	// 	log.Errorf("AWS initialization failed: %v", err)
	// 	os.Exit(1)
	// }

	// Initialize Router
	appRouter, err := InitAppRouter(dbConn)
	if err != nil {
		log.Errorf("AppRouter initialization failed: %v", err)
		os.Exit(1)
	}

	// Initialize Server
	server := InitServer(appRouter)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-quit
		log.Infof("Received signal '%v'. Shutting down server...", sig)
		if err := server.Shutdown(); err != nil {
			log.Warnf("Server forced to shutdown: %v", err)
		}
		log.Info("Server exiting")
	}()

	// Start the server
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		log.Errorf("Failed to run the server: %v", err)
		os.Exit(1)
	}
}
