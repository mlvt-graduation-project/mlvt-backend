package migration

import (
	"database/sql"
	"fmt"

	"mlvt/internal/infra/zap-logging/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDB applies all pending migrations.
func MigrateDB(db *sql.DB) error {
	// Initialize the migration driver
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrationPath := "file://../../migration"

	// Debug: Print the migration path
	fmt.Printf("Migration Path: %s\n", migrationPath)

	// Create a new migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Info("No new migrations to apply.")
	} else {
		log.Info("Migrations applied successfully.")
	}

	return nil
}
