package repo

import (
	"mlvt/internal/repo/admin_repo"
	"mlvt/internal/repo/media_repo"
	"mlvt/internal/repo/ping_repo"
	"mlvt/internal/repo/progress_repo"
	"mlvt/internal/repo/traffic_repo"
	"mlvt/internal/repo/user_repo"
	"mlvt/internal/repo/voucher_repo"
	"mlvt/internal/repo/wallet_repo"

	"github.com/google/wire"
)

// ProviderSetRepository is providers.
var ProviderSetRepository = wire.NewSet(
	user_repo.NewUserRepo,
	media_repo.NewMediaRepo,
	progress_repo.NewProgressRepo,
	admin_repo.NewAminRepo,
	traffic_repo.NewTrafficRepo,
	wallet_repo.NewWalletRepo,
	voucher_repo.NewVoucherRepo,
	ping_repo.NewPingRepo,
	// wire.Bind(new(UserRepository), new(*userRepo)),
	// wire.Bind(new(VideoRepository), new(*videoRepo)),
	// wire.Bind(new(AudioRepository), new(*audioRepo)),
	// wire.Bind(new(TranscriptionRepository), new(*transcriptionRepo)),
)
