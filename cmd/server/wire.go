// wire.go
//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	handler "mlvt/internal/handler/rest/v1"
	"mlvt/internal/infra/aws"
	"mlvt/internal/pkg/middleware"
	"mlvt/internal/repo"
	"mlvt/internal/router"
	"mlvt/internal/service"

	"github.com/google/wire"
)

func InitializeApp(db *sql.DB, s3Client *aws.S3Client) (*router.AppRouter, error) {
	wire.Build(
		repo.ProviderSetRepository,
		service.ProviderSetService,
		handler.ProviderSetHandler,
		middleware.ProviderSetMiddleware,
		router.ProviderSetRouter,
	)
	return &router.AppRouter{}, nil
}
