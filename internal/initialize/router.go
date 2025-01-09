package initialize

import (
	"database/sql"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/router"
)

// InitRouter sets up the application router using dependency injection.
func InitAppRouter(dbConn *sql.DB, mongoConn *mongodb.MongoDBClient) (*router.AppRouter, error) {
	appRouter, err := InitializeApp(dbConn, mongoConn)
	if err != nil {
		return nil, err
	}
	return appRouter, nil
}
