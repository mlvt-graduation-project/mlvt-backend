package cleanup

import (
	"log"
	"mlvt/internal/infra/aws"
	"mlvt/internal/infra/db"
	"mlvt/internal/infra/seeder"
	"mlvt/internal/repo"
)

func main() {

	// Initialize the database
	dbConn, err := db.InitializeDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize repositories
	userRepo := repo.NewUserRepo(dbConn)
	videoRepo := repo.NewVideoRepo(dbConn)

	// Initialize AWS S3 client
	s3Client, err := aws.NewS3Client()
	if err != nil {
		log.Fatalf("Failed to initialize AWS S3 client: %v", err)
	}

	// Initialize the seeder (can be reused for cleanup)
	userVideoSeeder := seeder.NewUserVideoSeeder(userRepo, videoRepo, s3Client)

	// Perform cleanup
	err = userVideoSeeder.CleanupSeededData()
	if err != nil {
		log.Fatalf("Failed to cleanup seeded data: %v", err)
	}

	log.Println("Cleanup of seeded data completed successfully.")
}
