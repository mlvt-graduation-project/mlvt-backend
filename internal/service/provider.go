package service

import "github.com/google/wire"

const SecretKey = "your-secret-key"

// ProviderSetService is providers.
var ProviderSetService = wire.NewSet(
	NewAuthService,
	NewUserService,
	NewVideoService,
	wire.Value(SecretKey),
)
