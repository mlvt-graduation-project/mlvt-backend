package initialize

import (
	"database/sql"
	"mlvt/internal/router"
)

// InitRouter sets up the application router using dependency injection.
func InitAppRouter(dbConn *sql.DB) (*router.AppRouter, error) {
	appRouter, err := InitializeApp(dbConn)
	if err != nil {
		return nil, err
	}
	return appRouter, nil
}
