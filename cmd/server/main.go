package main

import (
	"flag"
	"fmt"
	"mlvt/cmd/migration"
	"mlvt/internal/infra/aws"
	"mlvt/internal/infra/db"
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/reason"
	"mlvt/internal/infra/server/http"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/infra/zap-logging/zap"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	if env.EnvConfig == nil {
		fmt.Errorf("EnvConfig not loaded")
	}

	// Read environment variables
	logLevel = env.EnvConfig.LogLevel
	logPath = env.EnvConfig.LogPath

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

	log.Info(reason.MigrationsApplied.Message())

	// Initialize AWS S3 client
	s3Client, err := aws.NewS3Client()
	if err != nil {
		log.Errorf("Failed to initialize AWS S3 client: %v", err)
		os.Exit(1)
	}

	appRouter, err := InitializeApp(dbConn, s3Client)
	if err != nil {
		log.Errorf("Failed to initialize app: %v", err)
		os.Exit(1)
	}

	// Create a new Gin router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Đặt nguồn bạn muốn cho phép
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Cho phép gửi thông tin xác thực như cookie
		MaxAge:           12 * time.Hour,
	}))
	api := r.Group("/api")
	appRouter.RegisterUserRoutes(api)
	appRouter.RegisterVideoRoutes(api)
	appRouter.RegisterTranscriptionRoutes(api)
	appRouter.RegisterSwaggerRoutes(r.Group("/"))

	// Create the http server
	addr := ":" + env.EnvConfig.ServerPort
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
