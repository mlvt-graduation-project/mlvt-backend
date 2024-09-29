package migration

import (
	"database/sql"
	"fmt"
	"mlvt/internal/infra/reason"
	"mlvt/internal/infra/zap-logging/log"
	"time"
)

// Migration defines a database migration.
type Migration struct {
	ID   int
	Name string
	SQL  string
}

// Migrate runs all pending migrations and inserts sample data.
func Migrate(db *sql.DB) error {
	// Ensure the migrations table exists
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	// Define migrations
	migrations := []Migration{
		{
			ID:   1,
			Name: "create_users_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					first_name TEXT NOT NULL,
					last_name TEXT NOT NULL,
					email TEXT NOT NULL UNIQUE,
					password TEXT NOT NULL,
					status INTEGER NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);`,
		},
		{
			ID:   2,
			Name: "create_videos_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS videos (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER NOT NULL,
					title TEXT NOT NULL,
					duration INTEGER NOT NULL,
					link TEXT NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
				);`,
		},
		// New migration to alter the videos table
		{
			ID:   3,
			Name: "alter_videos_table_add_file_name_and_image",
			SQL: `
				ALTER TABLE videos
				ADD COLUMN file_name TEXT NOT NULL DEFAULT '',
				ADD COLUMN image TEXT NOT NULL DEFAULT '';
			`,
		},
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if err := applyMigration(db, migration); err != nil {
			return err
		}
	}

	// Insert sample data
	if err := insertSampleData(db); err != nil {
		return err
	}

	return nil
}

// insertSampleData inserts sample data into the database for testing purposes.
func insertSampleData(db *sql.DB) error {
	// fmt.Println("Inserting sample video data...")
	baseTitles := []string{
		"Introduction to Go", "Learning AWS S3 Basics", "Cooking Masterclass",
		"Travel Vlog: Paris", "Yoga for Beginners", "Machine Learning 101",
		"Advanced Golang Techniques", "Docker Basics", "ReactJS Tutorial",
		"Kubernetes Deployment", "Cloud Computing", "Database Optimization",
		"Web Development with Flask", "Microservices Architecture", "Python for Data Science",
		"Getting Started with TypeScript", "Mastering GraphQL", "DevOps Best Practices",
		"Building APIs with Node.js", "Serverless Computing Overview",
	}
	userIDs := []uint64{1, 2, 3, 4, 5}
	durationBase := 300 // Base duration in seconds

	// Generate 50 sample video entries
	for i := 1; i <= 50; i++ {
		title := fmt.Sprintf("%s - Part %d", baseTitles[i%len(baseTitles)], i)
		userID := userIDs[i%len(userIDs)]
		duration := durationBase + (i * 10)
		link := fmt.Sprintf("https://example.com/video%d.mp4", i)
		fileName := fmt.Sprintf("video%d.mp4", i)
		image := fmt.Sprintf("https://example.com/video%d-thumbnail.jpg", i)

		query := `INSERT INTO videos (user_id, title, duration, link, file_name, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query, userID, title, duration, link, fileName, image, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf(reason.InsertSampleFailed.Message()+" '%s': %v", title, err)
		}
		// fmt.Printf("Inserted sample video: %s\n", title)
	}
	log.Info(reason.InsertSampleDataSuccess.Message())

	return nil
}

// ensureMigrationsTable ensures that the migrations table exists.
func ensureMigrationsTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(createTableSQL)
	return err
}

// applyMigration applies a single migration if it hasn't been applied yet.
func applyMigration(db *sql.DB, migration Migration) error {
	// Check if the migration has already been applied
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", migration.Name).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// Migration already applied
		log.Infof(reason.Migration.Message(), " ", migration.Name, " ", reason.AlreadyApplied.Message())
		return nil
	}

	// Apply the migration
	log.Infof(reason.ApplyingMigration.Message(), " ", migration.Name)
	_, err = db.Exec(migration.SQL)
	if err != nil {
		return err
	}

	// Record the migration as applied
	_, err = db.Exec("INSERT INTO migrations (name) VALUES (?)", migration.Name)
	return err
}
