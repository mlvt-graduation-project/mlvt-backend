package repo

import (
	"mlvt/internal/repo/audio_repo"
	"mlvt/internal/repo/payment_repo/momo_repo"
	"mlvt/internal/repo/ping_repo"
	"mlvt/internal/repo/transcription_repo"
	"mlvt/internal/repo/user_repo"
	"mlvt/internal/repo/video_repo"

	"github.com/google/wire"
)

// ProviderSetRepository is providers.
var ProviderSetRepository = wire.NewSet(
	user_repo.NewUserRepo,
	video_repo.NewVideoRepo,
	audio_repo.NewAudioRepository,
	transcription_repo.NewTranscriptionRepository,
	momo_repo.NewMoMoRepo,
	ping_repo.NewPingRepo,
	// wire.Bind(new(UserRepository), new(*userRepo)),
	// wire.Bind(new(VideoRepository), new(*videoRepo)),
	// wire.Bind(new(AudioRepository), new(*audioRepo)),
	// wire.Bind(new(TranscriptionRepository), new(*transcriptionRepo)),
)
