package handler

import (
	"mlvt/internal/handler/rest/v1/admin_handler"
	"mlvt/internal/handler/rest/v1/media_handler"
	"mlvt/internal/handler/rest/v1/mlvt_handler"
	"mlvt/internal/handler/rest/v1/payment_handler"
	"mlvt/internal/handler/rest/v1/ping_handler"
	"mlvt/internal/handler/rest/v1/progress_handler"
	"mlvt/internal/handler/rest/v1/token_claim_handler"
	"mlvt/internal/handler/rest/v1/user_handler"
	"mlvt/internal/handler/rest/v1/voucher_handler"
	"mlvt/internal/handler/rest/v1/wallet_handler"

	"github.com/google/wire"
)

// ProviderSetHandler is Handler providers.
var ProviderSetHandler = wire.NewSet(
	user_handler.NewUserController,
	media_handler.NewMediaController,
	payment_handler.NewPaymentController,
	ping_handler.NewPingController,
	mlvt_handler.NewMlvtController,
	progress_handler.NewProgressService,
	admin_handler.NewAdminController,
	wallet_handler.NewWalletController,
	voucher_handler.NewVoucherController,
	token_claim_handler.New,
)
