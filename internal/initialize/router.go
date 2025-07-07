package initialize

import (
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/router"

	"github.com/jmoiron/sqlx"
)

// InitRouter sets up the application router using dependency injection.
func InitAppRouter(dbConn *sqlx.DB, mongoConn *mongodb.MongoDBClient) (*router.AppRouter, error) {
	appRouter, err := InitializeApp(dbConn, mongoConn)
	if err != nil {
		return nil, err
	}
	return appRouter, nil
}
