package handler

import (
	"mlvt/internal/handler/rest/v1/admin_handler"
	"mlvt/internal/handler/rest/v1/audio_handler"
	"mlvt/internal/handler/rest/v1/mlvt_handler"
	"mlvt/internal/handler/rest/v1/payment_handler/momo_handler"
	"mlvt/internal/handler/rest/v1/ping_handler"
	"mlvt/internal/handler/rest/v1/progress_handler"
	"mlvt/internal/handler/rest/v1/transcription_handler"
	"mlvt/internal/handler/rest/v1/user_handler"
	"mlvt/internal/handler/rest/v1/video_handler"
	"mlvt/internal/handler/rest/v1/wallet_handler"

	"github.com/google/wire"
)

// ProviderSetHandler is Handler providers.
var ProviderSetHandler = wire.NewSet(
	user_handler.NewUserController,
	video_handler.NewVideoController,
	audio_handler.NewAudioController,
	transcription_handler.NewTranscriptionController,
	ping_handler.NewPingController,
	mlvt_handler.NewMlvtController,
	progress_handler.NewProgressService,
	admin_handler.NewAdminController,
	wallet_handler.NewWalletController,
	momo_handler.NewMoMoPaymentHandler,
)
