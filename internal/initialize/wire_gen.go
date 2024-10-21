// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package initialize

import (
	"database/sql"
	"mlvt/internal/handler/rest/v1"
	"mlvt/internal/infra/aws"
	"mlvt/internal/pkg/middleware"
	"mlvt/internal/repo"
	"mlvt/internal/router"
	"mlvt/internal/service"
)

// Injectors from wire.go:

func InitializeApp(db *sql.DB) (*router.AppRouter, error) {
	userRepository := repo.NewUserRepo(db)
	s3ClientInterface, err := aws.NewS3Client()
	if err != nil {
		return nil, err
	}
	string2 := _wireStringValue
	authServiceInterface := service.NewAuthService(userRepository, string2)
	userService := service.NewUserService(userRepository, s3ClientInterface, authServiceInterface)
	userController := handler.NewUserController(userService)
	videoRepository := repo.NewVideoRepo(db)
	videoService := service.NewVideoService(videoRepository, s3ClientInterface)
	videoController := handler.NewVideoController(videoService)
	audioRepository := repo.NewAudioRepository(db)
	audioService := service.NewAudioService(audioRepository, s3ClientInterface)
	audioController := handler.NewAudioController(audioService)
	transcriptionRepository := repo.NewTranscriptionRepository(db)
	transcriptionService := service.NewTranscriptionService(transcriptionRepository, s3ClientInterface)
	transcriptionController := handler.NewTranscriptionController(transcriptionService)
	authUserMiddleware := middleware.NewAuthUserMiddleware(authServiceInterface)
	moMoRepo := repo.NewMoMoRepo()
	moMoPaymentService := service.NewMoMoPaymentService(moMoRepo)
	moMoPaymentController := handler.NewMoMoPaymentHandler(moMoPaymentService)
	swaggerRouter := router.NewSwaggerRouter()
	appRouter := router.NewAppRouter(userController, videoController, audioController, transcriptionController, authUserMiddleware, moMoPaymentController, swaggerRouter)
	return appRouter, nil
}

var (
	_wireStringValue = service.SecretKey
)
