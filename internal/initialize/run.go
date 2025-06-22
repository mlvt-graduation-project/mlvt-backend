package initialize

import (
	"fmt"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/service/notify_service"
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

	// Initialize notification service
	notifyService := notify_service.NewNotifyService()

	// Initialize Database
	dbConn, mongoConn, err := InitDatabase()
	if err != nil {
		log.Errorf("Database initialization failed: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Warnf("Error closing database connection: %v", err)
		}
	}()
	defer func() {
		if err := mongoConn.Close(); err != nil {
			log.Warnf("Error closing MongoDB connection: %v", err)
		}
	}()

	// Initialize AWS Clients
	// s3Client, err := InitAWS()
	// if err != nil {
	// 	log.Errorf("AWS initialization failed: %v", err)
	// 	os.Exit(1)
	// }

	// Initialize Router
	appRouter, err := InitAppRouter(dbConn, mongoConn)
	if err != nil {
		log.Errorf("AppRouter initialization failed: %v", err)
		os.Exit(1)
	}

	// Initialize Server
	server := InitServer(appRouter)

	// Send startup notification
	startupMessage := "🚀 <b>MLVT-Backend Server Starting</b>\n\n" +
		"✅ Logger initialized\n" +
		"✅ Database connections established\n" +
		"✅ Routes configured\n" +
		"✅ Server ready to accept connections\n" +
		"🌐 MLVT-Backend is now online"

	if err := notifyService.SendFormattedNotification(startupMessage); err != nil {
		log.Warnf("Failed to send startup notification: %v", err)
	} else {
		log.Info("Startup notification sent successfully")
	}

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-quit
		log.Infof("Received signal '%v'. Shutting down server...", sig)

		// Send shutdown notification
		shutdownMessage := fmt.Sprintf("🛑 <b>MLVT-Backend Server Shutting Down</b>\n\n"+
			"📡 Signal received: %v\n"+
			"⏳ Initiating graceful shutdown...\n"+
			"🔄 MLVT-Backend going offline", sig)

		if err := notifyService.SendFormattedNotification(shutdownMessage); err != nil {
			log.Warnf("Failed to send shutdown notification: %v", err)
		} else {
			log.Info("Shutdown notification sent successfully")
		}

		if err := server.Shutdown(); err != nil {
			log.Warnf("Server forced to shutdown: %v", err)
			// Send error notification
			errorMessage := fmt.Sprintf("❌ <b>MLVT-Backend Shutdown Error</b>\n\n"+
				"🔥 Error during shutdown: %v\n"+
				"⚠️ MLVT-Backend forced shutdown", err)
			if notifyErr := notifyService.SendFormattedNotification(errorMessage); notifyErr != nil {
				log.Warnf("Failed to send error notification: %v", notifyErr)
			}
		} else {
			// Send successful shutdown notification
			successMessage := "✅ <b>MLVT-Backend Shutdown Complete</b>\n\n" +
				"🔒 All connections closed gracefully\n" +
				"📊 Server exited successfully\n" +
				"💤 MLVT-Backend is now offline"
			if err := notifyService.SendFormattedNotification(successMessage); err != nil {
				log.Warnf("Failed to send success notification: %v", err)
			}
		}
		log.Info("Server exiting")
	}()

	// Start the server
	log.Info("Starting HTTP server...")
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		log.Errorf("Failed to run the server: %v", err)

		// Send server start error notification
		errorMessage := fmt.Sprintf("❌ <b>MLVT-Backend Start Failed</b>\n\n"+
			"🔥 Error: %v\n"+
			"⚠️ MLVT-Backend failed to start", err)
		if notifyErr := notifyService.SendFormattedNotification(errorMessage); notifyErr != nil {
			log.Warnf("Failed to send server start error notification: %v", notifyErr)
		}

		os.Exit(1)
	}
}
