package routes

import (
	"database/sql"
	"log"
	handlers_rest1 "mlvt/internal/handlers/rest/v1"
	"mlvt/internal/repository"
	"mlvt/internal/service"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, awsClient *s3.Client) *gin.Engine {
	//func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	trustedProxies := []string{"192.168.0.1", "10.0.0.1"}
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatal("Failed to set trusted proxies: %v", err)
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/test/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "Hello %s", name)
	})

	// Repositories
	//userRepo := repository.NewPostgresUserRepository(db)
	videoRepo := repository.NewPostgresVideoRepository(db)
	awsRepo := repository.NewAWSService(awsClient)

	// Services
	//authService := service.NewAuthService(userRepo, authConfig)
	//userService := service.NewUserService(userRepo)
	videoService := service.NewVideoService(videoRepo, awsRepo)

	// Handlers
	//userHandler := handlers.NewUserHandler(userService)
	//authHandler := handlers.NewAuthHandler(authService)
	videoHandler := handlers_rest1.NewVideoHandler(videoService)

	// User routes
	//router.POST("/users", userHandler.RegisterUser)
	//router.POST("/login", authHandler.Login)
	//router.POST("/logout", authHandler.Logout)

	// Video routes
	authGroup := router.Group("/auth")
	//authGroup.Use(auth.AuthMiddleware(userRepo))
	{
		authGroup.POST("/videos", videoHandler.UploadVideos)
		authGroup.GET("/videos/user/:userID", videoHandler.GetUserVideos)
	}

	return router

}
