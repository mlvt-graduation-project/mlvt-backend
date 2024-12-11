package migration

import (
	"database/sql"
	"fmt"

	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDB applies all pending migrations.
// It automatically handles dirty states by forcing the current version and retrying the migration.
func MigrateDB(db *sql.DB) error {
	// Initialize the migration driver
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrationPath := env.EnvConfig.MigrationsPath

	// Debug: Print the migration path
	log.Infof("Migration Path: %s", migrationPath)

	// Create a new migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() {
		// Close the migrate instance to clean up resources
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Errorf("Error closing migrate source: %v", srcErr)
		}
		if dbErr != nil {
			log.Errorf("Error closing migrate database: %v", dbErr)
		}
	}()

	// Attempt to apply migrations
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Info("No new migrations to apply.")
			return nil
		}

		// Check if the error is a dirty state
		// Since migrate.ErrDirty is not a pointer, use type assertion differently
		if migrateErr, ok := err.(migrate.ErrDirty); ok {
			log.Warnf("Database is dirty at version %d. Attempting to force the version.", migrateErr.Version)

			// Force the migration to the dirty version
			if forceErr := m.Force(int(migrateErr.Version)); forceErr != nil {
				return fmt.Errorf("failed to force migration to version %d: %w", migrateErr.Version, forceErr)
			}

			log.Infof("Forced migration to version %d successfully.", migrateErr.Version)

			// Retry the migration after forcing
			if retryErr := m.Up(); retryErr != nil {
				if retryErr == migrate.ErrNoChange {
					log.Info("Migrations applied successfully after forcing the version.")
					return nil
				}
				// If it's still dirty or another error occurred
				return fmt.Errorf("migration failed after forcing version: %w", retryErr)
			}

			log.Info("Migrations applied successfully after forcing the version.")
			return nil
		}

		// For other types of migration errors, log and return
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Info("Migrations applied successfully.")
	return nil
}
