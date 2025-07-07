package db

import (
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func InitializeDB() (*sqlx.DB, error) {
	driver := env.EnvConfig.DBDriver

	// Format the PostgreSQL connection string
	dsn := env.EnvConfig.DBConnection

	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Info("Database connection established successfully!")
	return db, nil
}
