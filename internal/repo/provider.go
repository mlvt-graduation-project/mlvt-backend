package repo

import "github.com/google/wire"

// ProviderSetRepository is peoviders.
var ProviderSetRepository = wire.NewSet(
	NewUserRepo,
	NewVideoRepo,
	wire.Bind(new(UserRepository), new(*UserRepo)),
	wire.Bind(new(VideoRepository), new(*VideoRepo)),
)
