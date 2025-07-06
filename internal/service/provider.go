package service

import (
	"mlvt/internal/infra/env"
	"mlvt/internal/service/admin_service"
	"mlvt/internal/service/auth_service"
	"mlvt/internal/service/media_service"
	"mlvt/internal/service/notify_service"
	"mlvt/internal/service/payment_service"
	"mlvt/internal/service/ping_service"
	"mlvt/internal/service/progress_service"
	"mlvt/internal/service/token_claim_service"
	"mlvt/internal/service/traffic_service"
	"mlvt/internal/service/user_service"
	"mlvt/internal/service/voucher_service"
	"mlvt/internal/service/wallet_service"

	"github.com/google/wire"
)

var SecretKey = env.EnvConfig.JWTSecret

// ProviderSetService is providers.
var ProviderSetService = wire.NewSet(
	auth_service.NewAuthService,
	user_service.NewUserService,
	media_service.NewMediaService,
	notify_service.NewNotifyService,
	payment_service.NewPaymentService,
	progress_service.NewProgressService,
	traffic_service.NewTrafficService,
	wallet_service.NewWalletService,
	voucher_service.NewVoucherService,
	admin_service.NewAminService,
	ping_service.NewPingService,
	token_claim_service.New,
	wire.Value(SecretKey),
)
