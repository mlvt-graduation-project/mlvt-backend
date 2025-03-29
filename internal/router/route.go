package router

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
	"mlvt/internal/handler/rest/v1/voucher_handler"
	"mlvt/internal/handler/rest/v1/wallet_handler"
	"mlvt/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	userController          *user_handler.UserController
	videoController         *video_handler.VideoController
	audioController         *audio_handler.AudioController
	transcriptionController *transcription_handler.TranscriptionController
	mlvtController          *mlvt_handler.MlvtController
	progressController      *progress_handler.ProgressController
	pingController          *ping_handler.PingController
	authMiddleware          *middleware.AuthUserMiddleware
	adminController         *admin_handler.AdminController
	walletController        *wallet_handler.WalletController
	voucherController       *voucher_handler.VoucherController
	momoPaymentController   *momo_handler.MoMoPaymentController
	swaggerRouter           *SwaggerRouter
}

func NewAppRouter(
	userController *user_handler.UserController,
	videoController *video_handler.VideoController,
	audioController *audio_handler.AudioController,
	transcriptionController *transcription_handler.TranscriptionController,
	mlvtController *mlvt_handler.MlvtController,
	progressController *progress_handler.ProgressController,
	pingController *ping_handler.PingController,
	authMiddleware *middleware.AuthUserMiddleware,
	adminController *admin_handler.AdminController,
	walletController *wallet_handler.WalletController,
	voucherController *voucher_handler.VoucherController,
	momoPaymentController *momo_handler.MoMoPaymentController,
	swaggerRouter *SwaggerRouter) *AppRouter {
	return &AppRouter{
		userController:          userController,
		videoController:         videoController,
		audioController:         audioController,
		transcriptionController: transcriptionController,
		mlvtController:          mlvtController,
		progressController:      progressController,
		pingController:          pingController,
		authMiddleware:          authMiddleware,
		adminController:         adminController,
		walletController:        walletController,
		voucherController:       voucherController,
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
		protected.GET("/:video_id/status", a.videoController.GetVideoStatus)                          // Get video status
		protected.PUT("/:video_id/status", a.videoController.UpdateVideoStatus)                       // Update video status
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
		protected.POST("/", a.transcriptionController.AddTranscription)                                         // Add a new transcription
		protected.GET("/:transcription_id", a.transcriptionController.GetTranscriptionByID)                     // Get transcription by ID
		protected.GET("/:transcription_id/user/:userID", a.transcriptionController.GetTranscriptionByUserID)    // Get transcription by transcription ID and user ID
		protected.GET("/:transcription_id/video/:videoID", a.transcriptionController.GetTranscriptionByVideoID) // Get transcription by transcription ID and video ID
		protected.GET("/user/:user_id", a.transcriptionController.ListTranscriptionsByUserID)                   // List transcriptions by user ID
		protected.GET("/video/:video_id", a.transcriptionController.ListTranscriptionsByVideoID)                // List transcriptions by video ID
		protected.DELETE("/:transcription_id", a.transcriptionController.DeleteTranscription)                   // Delete transcription by ID
		protected.POST("/generate-upload-url", a.transcriptionController.GenerateUploadURL)                     // Generate presigned upload URL
		protected.GET("/:transcription_id/download-url", a.transcriptionController.GenerateDownloadURL)         // Generate presigned download URL
		protected.PUT("/:transcription_id/status", a.transcriptionController.UpdateTranscriptionStatus)
	}
}

// RegisterAudioRoutes sets up the routes for audio-related operations
func (a *AppRouter) RegisterAudioRoutes(r *gin.RouterGroup) {
	protected := r.Group("/audios")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.POST("/", a.audioController.AddAudio)                                  // Add a new audio
		protected.GET("/:audio_id", a.audioController.GetAudio)                          // Get a specific audio by ID
		protected.DELETE("/:audio_id", a.audioController.DeleteAudio)                    // Delete an audio
		protected.GET("/user/:user_id", a.audioController.ListAudiosByUserID)            // Get all audios by user
		protected.GET("/video/:video_id", a.audioController.ListAudiosByVideoID)         // Get all audios by video
		protected.GET("/:audio_id/user/:user_id", a.audioController.GetAudioByUser)      // Get specific audio by audio ID and user ID
		protected.GET("/:audio_id/video/:video_id", a.audioController.GetAudioByVideoID) // Get specific audio by audio ID and video ID
		protected.POST("/generate-presigned-url", a.audioController.GenerateUploadURL)   // Generate presigned URL for audio upload
		protected.GET("/:audio_id/download-url", a.audioController.GenerateDownloadURL)  // Generate presigned URL for audio download
	}
}

func (a *AppRouter) RegisterPingStatusRoutes(r *gin.RouterGroup) {
	public := r.Group("/ping")
	{
		public.GET("/speech-to-text/:id", a.pingController.PingSpeechToText)
		public.GET("/text-to-text/:id", a.pingController.PingTextToText)
		public.GET("/text-to-speech/:id", a.pingController.PingTextToSpeech)
		public.GET("/voice-cloning/:id", a.pingController.PingVoiceCloning)
		public.GET("/lipsync/:id", a.pingController.PingLipSync)
		public.GET("/full-pipeline/:id", a.pingController.PingFullPipeline)
	}
}

func (a *AppRouter) RegiserMlvtRoutes(r *gin.RouterGroup) {
	public := r.Group("/mlvt")
	{
		public.POST("/ttt/:transcription_id", a.mlvtController.ProcessTextToText)
		public.POST("/stt/:video_id", a.mlvtController.ProcessSpeechToText)
		public.POST("/tts/:transcription_id", a.mlvtController.ProcessTextToSpeech)
		public.POST("/lipsync/:video_id/:audio_id", a.mlvtController.ProcessLipSync)
		public.POST("/pipeline/full/:video_id", a.mlvtController.ProcessFullPipeline)
	}
}

// RegisterProgressRoutes sets up the routes for all progress-related operations
func (a *AppRouter) RegisterProgressRoutes(r *gin.RouterGroup) {
	public := r.Group("/progress")
	{
		public.POST("/:user_id", a.progressController.GetUserProgress)
	}
}

// RegisterAdminRoutes sets up the routes for all permission-related operations
func (a *AppRouter) RegisterAdminRoutes(r *gin.RouterGroup) {
	protected := r.Group("/admin")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.GET("/:adminID/config", a.adminController.GetServerConfig)
		protected.POST("/:adminID/config", a.adminController.UpdateServerConfig)
		protected.GET("/:adminID/models", a.adminController.GetModelList)
		protected.POST("/:adminID/models", a.adminController.AddModelOption)
		protected.PUT("/:adminID/models/:modelOptionID", a.adminController.UpdateModelOption)
	}
}

func (a *AppRouter) RegisterWalletRoutes(r *gin.RouterGroup) {
	protected := r.Group("/wallet")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.GET("/deposit", a.walletController.Deposit)
		protected.POST("/use-token", a.walletController.UseToken)
		protected.GET("/balance", a.walletController.GetBalance)
	}
}

func (a *AppRouter) RegisteVoucherRoutes(r *gin.RouterGroup) {
	protected := r.Group("/voucher")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.POST("/voucher", a.voucherController.CreateVoucher)
		protected.POST("/voucher/use/:code", a.voucherController.UseVoucher)
		protected.PATCH("/voucher/:id", a.voucherController.UpdateVoucher)
		protected.GET("/voucher", a.voucherController.GetAllVouchers)
		protected.GET("/voucher/:id", a.voucherController.GetVoucherByID)
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
