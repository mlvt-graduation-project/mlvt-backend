// wire.go
//go:build wireinject
// +build wireinject

package initialize

import (
	"database/sql"
	handler "mlvt/internal/handler/rest/v1"
	"mlvt/internal/infra/aws"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/pkg/middleware"
	"mlvt/internal/repo"
	"mlvt/internal/router"
	"mlvt/internal/service"

	"github.com/google/wire"
)

func InitializeApp(db *sql.DB, mongoConn *mongodb.MongoDBClient) (*router.AppRouter, error) {
	wire.Build(
		aws.ProviderSetAwsBucket,
		repo.ProviderSetRepository,
		service.ProviderSetService,
		handler.ProviderSetHandler,
		middleware.ProviderSetMiddleware,
		router.ProviderSetRouter,
	)
	return &router.AppRouter{}, nil
}
