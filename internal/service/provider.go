package service

import (
	"mlvt/internal/infra/env"

	"github.com/google/wire"
)

var SecretKey = env.EnvConfig.JWTSecret

// ProviderSetService is providers.
var ProviderSetService = wire.NewSet(
	NewAuthService,
	NewUserService,
	NewVideoService,
	wire.Value(SecretKey),
)
