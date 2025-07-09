package initialize

import (
	"fmt"
	"mlvt/internal/infra/db"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/infra/env"

	"github.com/jmoiron/sqlx"
)

// InitDatabase establishes a database connection and runs migrations.
func InitDatabase() (*sqlx.DB, *mongodb.MongoDBClient, error) {
	dbConn, err := db.InitializeDB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize the database: %w", err)
	}

	mongoConn := mongodb.NewMongoDBClient(env.EnvConfig.MongoDBEndPoint)

	// Run migrations
	// if err := migration.MigrateDB(dbConn); err != nil {
	// 	log.Errorf("Migration failed: %v", err)
	// 	return nil, fmt.Errorf("migration failed: %w", err)
	// }

	// log.Info("Migrations applied successfully.")

	return dbConn, mongoConn, nil
}
