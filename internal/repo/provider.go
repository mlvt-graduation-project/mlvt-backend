package repo

import (
	"mlvt/internal/repo/admin_repo"
	"mlvt/internal/repo/media_repo"
	"mlvt/internal/repo/payment_repo"
	"mlvt/internal/repo/ping_repo"
	"mlvt/internal/repo/progress_repo"
	"mlvt/internal/repo/token_claim_repo"
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
	payment_repo.NewPaymentRepo,
	progress_repo.NewProgressRepo,
	traffic_repo.NewTrafficRepo,
	wallet_repo.NewWalletRepo,
	voucher_repo.NewVoucherRepo,
	admin_repo.NewAminRepo,
	ping_repo.NewPingRepo,
	token_claim_repo.New,
)
