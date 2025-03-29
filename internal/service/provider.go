package service

import (
	"mlvt/internal/infra/env"
	"mlvt/internal/service/admin_service"
	"mlvt/internal/service/audio_service"
	"mlvt/internal/service/auth_service"
	"mlvt/internal/service/payment_service/momo_service"
	"mlvt/internal/service/ping_service"
	"mlvt/internal/service/progress_service"
	"mlvt/internal/service/traffic_service"
	"mlvt/internal/service/transcription_service"
	"mlvt/internal/service/user_service"
	"mlvt/internal/service/video_service"
	"mlvt/internal/service/voucher_service"
	"mlvt/internal/service/wallet_service"

	"github.com/google/wire"
)

var SecretKey = env.EnvConfig.JWTSecret

// ProviderSetService is providers.
var ProviderSetService = wire.NewSet(
	auth_service.NewAuthService,
	user_service.NewUserService,
	video_service.NewVideoService,
	audio_service.NewAudioService,
	transcription_service.NewTranscriptionService,
	progress_service.NewProgressService,
	traffic_service.NewTrafficService,
	wallet_service.NewWalletService,
	voucher_service.NewVoucherService,
	admin_service.NewAminService,
	momo_service.NewMoMoPaymentService,
	ping_service.NewPingService,
	wire.Value(SecretKey),
)
