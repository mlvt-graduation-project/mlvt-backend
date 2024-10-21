package initialize

import (
	"database/sql"
	"fmt"
	"mlvt/cmd/migration"
	"mlvt/internal/infra/db"
	"mlvt/internal/infra/zap-logging/log"
)

// InitDatabase establishes a database connection and runs migrations.
func InitDatabase() (*sql.DB, error) {
	dbConn, err := db.InitializeDB()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the database: %w", err)
	}

	// Run migrations
	if err := migration.MigrateDB(dbConn); err != nil {
		log.Errorf("Migration failed: %v", err)
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	log.Info("Migrations applied successfully.")

	return dbConn, nil
}
