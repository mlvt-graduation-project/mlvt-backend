package db

import (
	"database/sql"
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"

	_ "github.com/mattn/go-sqlite3"
)

// InitializeDB initializes the SQLite3 database and returns a database connection.
func InitializeDB() (*sql.DB, error) {
	dbPath := env.EnvConfig.DBConnection
	dbDriver := env.EnvConfig.DBDriver

	log.Infof("DBConnection: %s, DBDriver: %s", dbPath, dbDriver)

	// Open a connection to the database file (creates the file if it doesn't exist)
	db, err := sql.Open(dbDriver, dbPath)
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
