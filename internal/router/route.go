package router

import (
	handler "mlvt/internal/handler/rest/v1"
	"mlvt/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	userController          *handler.UserController
	videoController         *handler.VideoController
	transcriptionController *handler.TranscriptionController
	authMiddleware          *middleware.AuthUserMiddleware
	swaggerRouter           *SwaggerRouter
}

func NewAppRouter(userController *handler.UserController, videoController *handler.VideoController, transcriptionController *handler.TranscriptionController, authMiddleware *middleware.AuthUserMiddleware, swaggerRouter *SwaggerRouter) *AppRouter {
	return &AppRouter{
		userController:          userController,
		videoController:         videoController,
		transcriptionController: transcriptionController,
		authMiddleware:          authMiddleware,
		swaggerRouter:           swaggerRouter,
	}
}

// RegisterUserRoutes sets up the routes for user-related operations
func (a *AppRouter) RegisterUserRoutes(r *gin.RouterGroup) {
	public := r.Group("/users")
	{
		public.POST("/register", a.userController.RegisterUser)
		public.POST("/login", a.userController.Login)
	}

	protected := r.Group("/users")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.GET("/:id", a.userController.GetUser)
		protected.PUT("/:id", a.userController.UpdateUser)
		protected.DELETE("/:id", a.userController.DeleteUser)
	}
}

// RegisterVideoRoutes sets up the routes for video-related operations
func (a *AppRouter) RegisterVideoRoutes(r *gin.RouterGroup) {
	protected := r.Group("/videos")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.POST("/", a.videoController.AddVideo)
		protected.PUT("/:id", a.videoController.UpdateVideo)
		protected.GET("/:id", a.videoController.GetVideo)
		protected.DELETE("/:id", a.videoController.DeleteVideo)
		protected.GET("/user/:userID", a.videoController.GetVideosByUser)
		protected.POST("/generate-presigned-url", a.videoController.GeneratePresignedURLHandler)
	}
}

// RegisterTranscriptionRoutes sets up the routes for transcriptions-related operations
func (a *AppRouter) RegisterTranscriptionRoutes(r *gin.RouterGroup) {
	protected := r.Group("/transcriptions")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.POST("/", a.transcriptionController.CreateTranscription)
		protected.POST("/generate-presigned-url/:videoID", a.transcriptionController.GeneratePresignedURL)
		protected.GET("/:userID/:transcriptionID", a.transcriptionController.GetTranscription)
		protected.GET("/by-user/:userID", a.transcriptionController.ListTranscriptionsByUser)
		protected.GET("/by-video/:videoID/:transcriptionID", a.transcriptionController.GetTranscriptionByVideo)
		protected.GET("/by-video/:videoID", a.transcriptionController.ListTranscriptionsByVideo)
		protected.DELETE("/:transcriptionID", a.transcriptionController.DeleteTranscription)
	}
}

// RegisterSwaggerRoutes sets up the route for Swagger API documentation
func (a *AppRouter) RegisterSwaggerRoutes(r *gin.RouterGroup) {
	// Check if SwaggerRouter is initialized before registering
	if a.swaggerRouter != nil {
		a.swaggerRouter.Register(r)
	}
}
