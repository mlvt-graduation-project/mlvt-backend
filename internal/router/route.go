package router

import (
	handler "mlvt/internal/handler/rest/v1"
	"mlvt/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	userController          *handler.UserController
	videoController         *handler.VideoController
	audioController         *handler.AudioController
	transcriptionController *handler.TranscriptionController
	authMiddleware          *middleware.AuthUserMiddleware
	momoPaymentController   *handler.MoMoPaymentController
	swaggerRouter           *SwaggerRouter
}

func NewAppRouter(userController *handler.UserController, videoController *handler.VideoController, audioController *handler.AudioController, transcriptionController *handler.TranscriptionController, authMiddleware *middleware.AuthUserMiddleware, momoPaymentController *handler.MoMoPaymentController, swaggerRouter *SwaggerRouter) *AppRouter {
	return &AppRouter{
		userController:          userController,
		videoController:         videoController,
		audioController:         audioController,
		transcriptionController: transcriptionController,
		authMiddleware:          authMiddleware,
		momoPaymentController:   momoPaymentController,
		swaggerRouter:           swaggerRouter,
	}
}

// RegisterUserRoutes sets up the routes for user-related operations
func (a *AppRouter) RegisterUserRoutes(r *gin.RouterGroup) {
	public := r.Group("/users")
	{
		public.POST("/register", a.userController.RegisterUser)
		public.POST("/login", a.userController.LoginUser)
	}

	protected := r.Group("/users")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.GET("/:user_id", a.userController.GetUser)
		protected.PUT("/:user_id", a.userController.UpdateUser)
		protected.DELETE("/:user_id", a.userController.DeleteUser)
		protected.PUT("/:user_id/change-password", a.userController.ChangePassword)
		protected.PUT("/:user_id/update-avatar", a.userController.UpdateAvatar)                    // Avatar upload (presigned URL)
		protected.GET("/:user_id/avatar-download-url", a.userController.GenerateAvatarDownloadURL) // Avatar download (presigned URL)
		protected.GET("/:user_id/avatar", a.userController.LoadAvatar)                             // Load avatar directly
	}
}

// RegisterVideoRoutes sets up the routes for video-related operations
func (a *AppRouter) RegisterVideoRoutes(r *gin.RouterGroup) {
	protected := r.Group("/videos")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.POST("/", a.videoController.AddVideo)                                               // Add a new video
		protected.GET("/:video_id", a.videoController.GetVideoByID)                                   // Get video by ID
		protected.GET("/user/:user_id", a.videoController.ListVideosByUserID)                         // List videos by user ID
		protected.DELETE("/:video_id", a.videoController.DeleteVideo)                                 // Delete video by ID
		protected.POST("/generate-upload-url/video", a.videoController.GenerateUploadURLForVideo)     // Generate presigned upload URL for video
		protected.POST("/generate-upload-url/image", a.videoController.GenerateUploadURLForImage)     // Generate presigned upload URL for image
		protected.GET("/:video_id/download-url/video", a.videoController.GenerateDownloadURLForVideo) // Generate presigned download URL for video
		protected.GET("/:video_id/download-url/image", a.videoController.GenerateDownloadURLForImage) // Generate presigned download URL for image
	}
}

// RegisterTranscriptionRoutes sets up the routes for transcription-related operations
func (a *AppRouter) RegisterTranscriptionRoutes(r *gin.RouterGroup) {
	protected := r.Group("/transcriptions")
	protected.Use(a.authMiddleware.MustAuth()) // Require authentication
	{
		protected.POST("/", a.transcriptionController.AddTranscription)                                        // Add a new transcription
		protected.GET("/:transcriptionID", a.transcriptionController.GetTranscriptionByID)                     // Get transcription by ID
		protected.GET("/:transcriptionID/user/:userID", a.transcriptionController.GetTranscriptionByUserID)    // Get transcription by transcription ID and user ID
		protected.GET("/:transcriptionID/video/:videoID", a.transcriptionController.GetTranscriptionByVideoID) // Get transcription by transcription ID and video ID
		protected.GET("/user/:user_id", a.transcriptionController.ListTranscriptionsByUserID)                  // List transcriptions by user ID
		protected.GET("/video/:video_id", a.transcriptionController.ListTranscriptionsByVideoID)               // List transcriptions by video ID
		protected.DELETE("/:transcriptionID", a.transcriptionController.DeleteTranscription)                   // Delete transcription by ID
		protected.POST("/generate-upload-url", a.transcriptionController.GenerateUploadURL)                    // Generate presigned upload URL
		protected.GET("/:transcriptionID/download-url", a.transcriptionController.GenerateDownloadURL)         // Generate presigned download URL
	}
}

// RegisterAudioRoutes sets up the routes for audio-related operations
func (a *AppRouter) RegisterAudioRoutes(r *gin.RouterGroup) {
	protected := r.Group("/audios")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.POST("/", a.audioController.AddAudio)                                // Add a new audio
		protected.GET("/:audioID", a.audioController.GetAudio)                         // Get a specific audio by ID
		protected.DELETE("/:audioID", a.audioController.DeleteAudio)                   // Delete an audio
		protected.GET("/user/:userID", a.audioController.ListAudiosByUserID)           // Get all audios by user
		protected.GET("/video/:videoID", a.audioController.ListAudiosByVideoID)        // Get all audios by video
		protected.GET("/:audioID/user/:userID", a.audioController.GetAudioByUser)      // Get specific audio by audio ID and user ID
		protected.GET("/:audioID/video/:videoID", a.audioController.GetAudioByVideo)   // Get specific audio by audio ID and video ID
		protected.POST("/generate-presigned-url", a.audioController.GenerateUploadURL) // Generate presigned URL for audio upload
		protected.GET("/:audioID/download-url", a.audioController.GenerateDownloadURL) // Generate presigned URL for audio download
	}
}

// RegisterPaymentRoutes sets up the routes for all payment-related operations
func (a *AppRouter) RegisterPaymentRoutes(r *gin.RouterGroup) {
	payment := r.Group("/payments")
	{
		// Group for MoMo-specific routes
		momo := payment.Group("/momo")
		{
			momo.POST("/create", a.momoPaymentController.CreateMoMoPayment)     // Create MoMo payment and return QR code
			momo.POST("/check-status", a.momoPaymentController.CheckMoMoStatus) // Check status of MoMo payment
			momo.POST("/refund", a.momoPaymentController.RefundMoMoPayment)     // Refund MoMo payment
		}

		// More payment methods can be added here...
	}
}

// RegisterSwaggerRoutes sets up the route for Swagger API documentation
func (a *AppRouter) RegisterSwaggerRoutes(r *gin.RouterGroup) {
	// Check if SwaggerRouter is initialized before registering
	if a.swaggerRouter != nil {
		a.swaggerRouter.Register(r)
	}
}
