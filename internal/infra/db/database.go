package db

import (
	"database/sql"
	"mlvt/internal/infra/zap-logging/log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitializeDB initializes the SQLite3 database and returns a database connection.
func InitializeDB() (*sql.DB, error) {
	// Determine the absolute path to the project root
	projectRoot, err := getProjectRoot()
	if err != nil {
		return nil, err
	}

	// Construct the path to the database file in the project root
	dbPath := filepath.Join(projectRoot, "mlvt.db")

	// Open a connection to the SQLite3 database file (creates the file if it doesn't exist)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Info("Database connection established successfully!")
	return db, nil
}

// getProjectRoot returns the absolute path to the project root directory.
func getProjectRoot() (string, error) {
	// Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Determine the project root based on the current working directory
	// Assuming the project root is the parent directory of the `cmd/server/`
	return filepath.Abs(filepath.Join(workingDir, "../../"))
}
