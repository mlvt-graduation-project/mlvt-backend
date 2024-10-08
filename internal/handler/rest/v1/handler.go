package handler

import "github.com/google/wire"

// ProviderSetHandler is Handler providers.
var ProviderSetHandler = wire.NewSet(
	NewUserController,
	NewVideoController,
	NewAudioController,
	NewTranscriptionController,
	NewMoMoPaymentHandler,
)
