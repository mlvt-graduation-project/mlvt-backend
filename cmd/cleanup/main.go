package main

import (
	"fmt"
	"mlvt/internal/infra/seeder"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/initialize"
	"mlvt/internal/repo"
	"os"
)

func main() {

	// Initialize Logger
	if err := initialize.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Logger initialization failed: %v\n", err)
		os.Exit(1)
	}

	// Initialize Database
	dbConn, err := initialize.InitDatabase()
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
	s3Client, err := initialize.InitAWS()
	if err != nil {
		log.Errorf("AWS initialization failed: %v", err)
		os.Exit(1)
	}

	userRepo := repo.NewUserRepo(dbConn)
	videoRepo := repo.NewVideoRepo(dbConn)

	// Initialize the seeder (can be reused for cleanup)
	userVideoSeeder := seeder.NewUserVideoSeeder(userRepo, videoRepo, s3Client)

	// Perform cleanup
	err = userVideoSeeder.CleanupSeededData()
	if err != nil {
		log.Errorf("Failed to cleanup seeded data: %v", err)
	}

	log.Infof("Cleanup of seeded data completed successfully.")
}
