package router

import (
	handler "mlvt/internal/handler/rest/v1"
	"mlvt/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// AppRouter holds all controllers for routing
type AppRouter struct {
	userController  *handler.UserController
	videoController *handler.VideoController
	authMiddleware  *middleware.AuthUserMiddleware
}

// NewAppRouter creates a new AppRouter instance
func NewAppRouter(userController *handler.UserController, videoController *handler.VideoController, authMiddleware *middleware.AuthUserMiddleware) *AppRouter {
	return &AppRouter{
		userController:  userController,
		videoController: videoController,
		authMiddleware:  authMiddleware,
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
