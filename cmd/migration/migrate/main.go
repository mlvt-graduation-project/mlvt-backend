package main

import (
	"database/sql"
	"fmt"
	"os"

	v1 "mlvt/cmd/migration/migrate/v1"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/initialize"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

func main() {
	// Initialize Logger
	if err := initialize.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Logger initialization failed: %v\n", err)
		os.Exit(1)
	}

	// Open SQLite database
	dsn := "mlvt.db"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Errorf("Failed to open database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Warnf("Error closing database: %v", err)
		}
	}()

	// Run the first version migration
	if err := v1.MigrateV1(db); err != nil {
		log.Errorf("MigrateV1 failed: %v", err)
		os.Exit(1)
	}

	log.Info("MigrateV1 completed successfully.")
}
