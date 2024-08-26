package main

import (
	"flag"
	"fmt"
	"mlvt/cmd/migration"
	"mlvt/internal/infra/db"
	"mlvt/internal/infra/server/http"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/infra/zap-logging/zap"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	// Name is the name of the project
	Name = "mlvt"
	// Version is the version of the project
	Version = "0.0.0"
	// confFlag is the config flag
	confFlag string
	// log level
	logLevel string
	// log path
	logPath string
)

func init() {
	flag.StringVar(&confFlag, "c", "../../configs/config.yaml", "config path, eg: -c config.yaml")
}

func main() {
	flag.Parse()

	// Load the .env file from the root directory
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}

	// Read environment variables
	logLevel = os.Getenv("LOG_LEVEL")
	logPath = os.Getenv("LOG_PATH")

	fmt.Println("LOG_PATH:", logPath) // Debug: print the log path

	// Ensure the log directory exists
	logDir := filepath.Dir(logPath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			fmt.Printf("Failed to create log directory: %v", err)
		}
	}

	// Initialize logging
	log.SetLogger(zap.NewLogger(
		log.ParseLevel(logLevel), zap.WithName(Name), zap.WithPath(logPath), zap.WithCallerFullPath()))

	dbConn, err := db.InitializeDB()
	if err != nil {
		log.Errorf("Failed to initialize the database: %v", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	if err := migration.Migrate(dbConn); err != nil {
		log.Errorf("Migration failed: %v", err)
	}

	log.Info("Migrations applied successfully!")

	appRouter, err := InitializeApp(dbConn)
	if err != nil {
		log.Errorf("Failed to initialize app: %v", err)
		os.Exit(1)
	}

	// Create a new Gin router
	r := gin.Default()
	api := r.Group("/api")
	appRouter.RegisterUserRoutes(api)
	appRouter.RegisterVideoRoutes(api)

	// Create the http server
	addr := ":8080" // Or read from an environment variable
	server := http.NewServer(r, addr)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Info("Shutting down server...")
		if err := server.Shutdown(); err != nil {
			log.Warnf("Server forced to shutdown: %v", err)
		}
		log.Info("Server exiting")
	}()

	// Start the server
	if err := server.Start(); err != nil {
		log.Errorf("Failed to run the server: %v", err)
	}
}
