package main

import (
	"fmt"
	"os"
	"path/filepath"

	"mlvt/internal/infra/env"
	"mlvt/internal/infra/seeder"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/initialize"
	"mlvt/internal/repo"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
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

	// Initialize repositories
	userRepo := repo.NewUserRepo(dbConn)
	videoRepo := repo.NewVideoRepo(dbConn)

	// Initialize the seeder
	userVideoSeeder := seeder.NewUserVideoSeeder(userRepo, videoRepo, s3Client)

	// Define the avatars and videos folders
	avatarsFolder := filepath.Join(env.EnvConfig.RootDir, "assets", "avatars")
	videosFolder := filepath.Join(env.EnvConfig.RootDir, "assets", "videos")

	// Check if avatars folder exists
	if _, err := os.Stat(avatarsFolder); os.IsNotExist(err) {
		log.Errorf("Avatars folder %s does not exist", avatarsFolder)
	}

	// Check if videos folder exists
	if _, err := os.Stat(videosFolder); os.IsNotExist(err) {
		log.Errorf("Videos folder %s does not exist", videosFolder)
	}

	// Seed users and videos
	err = userVideoSeeder.SeedUsersAndVideosFromFolders(avatarsFolder, videosFolder)
	if err != nil {
		log.Errorf("Failed to seed users and videos: %v", err)
	}

	log.Info("User and video seeding completed successfully.")
}
