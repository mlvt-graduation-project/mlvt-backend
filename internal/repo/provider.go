package repo

import "github.com/google/wire"

// ProviderSetRepository is providers.
var ProviderSetRepository = wire.NewSet(
	NewUserRepo,
	NewVideoRepo,
	NewAudioRepository,
	NewTranscriptionRepository,
	NewMoMoRepo,
	// wire.Bind(new(UserRepository), new(*userRepo)),
	// wire.Bind(new(VideoRepository), new(*videoRepo)),
	// wire.Bind(new(AudioRepository), new(*audioRepo)),
	// wire.Bind(new(TranscriptionRepository), new(*transcriptionRepo)),
)
