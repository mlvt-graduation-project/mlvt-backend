package service

import "github.com/google/wire"

const SecretKey = "veryvery"

// ProviderSetService is providers.
var ProviderSetService = wire.NewSet(
	NewAuthService,
	NewUserService,
	NewVideoService,
	wire.Value(SecretKey),
)
