package repo

import "github.com/google/wire"

// ProviderSetRepository is providers.
var ProviderSetRepository = wire.NewSet(
	NewUserRepo,
	NewVideoRepo,
	NewTranscriptionRepository,
	wire.Bind(new(UserRepository), new(*UserRepo)),
	wire.Bind(new(VideoRepository), new(*VideoRepo)),
	wire.Bind(new(TranscriptionRepository), new(*transcriptionRepo)),
)
