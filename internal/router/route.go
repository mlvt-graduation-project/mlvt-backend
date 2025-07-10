package router

import (
	"mlvt/internal/handler/rest/v1/admin_handler"
	"mlvt/internal/handler/rest/v1/media_handler"
	"mlvt/internal/handler/rest/v1/mlvt_handler"
	"mlvt/internal/handler/rest/v1/payment_handler"
	"mlvt/internal/handler/rest/v1/ping_handler"
	"mlvt/internal/handler/rest/v1/process_handler"
	"mlvt/internal/handler/rest/v1/progress_handler"
	"mlvt/internal/handler/rest/v1/token_claim_handler"
	"mlvt/internal/handler/rest/v1/user_handler"
	"mlvt/internal/handler/rest/v1/voucher_handler"
	"mlvt/internal/handler/rest/v1/wallet_handler"
	"mlvt/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	userController     *user_handler.UserController
	mediaController    *media_handler.MediaController
	mlvtController     *mlvt_handler.MlvtController
	progressController *progress_handler.ProgressController
	processController  *process_handler.ProcessController
	pingController     *ping_handler.PingController
	authMiddleware     *middleware.AuthUserMiddleware
	adminController    *admin_handler.AdminController
	walletController   *wallet_handler.WalletController
	voucherController  *voucher_handler.VoucherController
	tokenController    *token_claim_handler.TokenController
	paymentController  *payment_handler.PaymentController

	swaggerRouter *SwaggerRouter
}

func NewAppRouter(
	userController *user_handler.UserController,
	mediaController *media_handler.MediaController,
	mlvtController *mlvt_handler.MlvtController,
	progressController *progress_handler.ProgressController,
	processController *process_handler.ProcessController,
	pingController *ping_handler.PingController,
	authMiddleware *middleware.AuthUserMiddleware,
	adminController *admin_handler.AdminController,
	walletController *wallet_handler.WalletController,
	voucherController *voucher_handler.VoucherController,
	tokenController *token_claim_handler.TokenController,
	paymentController *payment_handler.PaymentController,
	swaggerRouter *SwaggerRouter) *AppRouter {
	return &AppRouter{
		userController:     userController,
		mediaController:    mediaController,
		mlvtController:     mlvtController,
		progressController: progressController,
		processController:  processController,
		pingController:     pingController,
		authMiddleware:     authMiddleware,
		adminController:    adminController,
		walletController:   walletController,
		voucherController:  voucherController,
		tokenController:    tokenController,
		paymentController:  paymentController,
		swaggerRouter:      swaggerRouter,
	}
}

// RegisterUserRoutes sets up the routes for user-related operations
func (a *AppRouter) RegisterUserRoutes(r *gin.RouterGroup) {
	public := r.Group("/users")
	{
		public.POST("/register", a.userController.RegisterUser)
		public.POST("/login", a.userController.LoginUser)
		public.POST("/verify-account", a.userController.VerifyAccountSignUp)
		public.POST("/resend-verification", a.userController.ResendValidationEmail)
	}

	protected := r.Group("/users")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.GET("", a.userController.GetAllUsers)
		protected.GET("/:user_id", a.userController.GetUser)
		protected.GET("/user-details", a.userController.GetUserDetails)
		protected.PUT("/:user_id", a.userController.UpdateUser)
		protected.DELETE("/:user_id", a.userController.DeleteUser)
		protected.PUT("/:user_id/change-password", a.userController.ChangePassword)
		protected.PUT("/:user_id/update-avatar", a.userController.UpdateAvatar)                    // Avatar upload (presigned URL)
		protected.GET("/:user_id/avatar-download-url", a.userController.GenerateAvatarDownloadURL) // Avatar download (presigned URL)
		protected.GET("/:user_id/avatar", a.userController.LoadAvatar)                             // Load avatar directly
	}
}

// RegisterMediaRoutes sets up the routes for media-related operations
func (a *AppRouter) RegisterMediaRoutes(r *gin.RouterGroup) {
	videoProtected := r.Group("/videos")
	videoProtected.Use(a.authMiddleware.MustAuth())
	{
		videoProtected.POST("/", a.mediaController.AddVideo)                                               // Add a new video
		videoProtected.GET("/:video_id", a.mediaController.GetVideoByID)                                   // Get video by ID
		videoProtected.GET("/user/:user_id", a.mediaController.ListVideosByUserID)                         // List videos by user ID
		videoProtected.DELETE("/:video_id", a.mediaController.DeleteVideo)                                 // Delete video by ID
		videoProtected.GET("/:video_id/status", a.mediaController.GetVideoStatus)                          // Get video status
		videoProtected.PUT("/:video_id/status", a.mediaController.UpdateVideoStatus)                       // Update video status
		videoProtected.POST("/generate-upload-url/video", a.mediaController.GenerateUploadURLForVideo)     // Generate presigned upload URL for video
		videoProtected.POST("/generate-upload-url/image", a.mediaController.GenerateUploadURLForImage)     // Generate presigned upload URL for image
		videoProtected.GET("/:video_id/download-url/video", a.mediaController.GenerateDownloadURLForVideo) // Generate presigned download URL for video
		videoProtected.GET("/:video_id/download-url/image", a.mediaController.GenerateDownloadURLForImage) // Generate presigned download URL for image
	}

	transcriptionProtected := r.Group("/transcriptions")
	transcriptionProtected.Use(a.authMiddleware.MustAuth()) // Require authentication
	{
		transcriptionProtected.POST("/", a.mediaController.AddTranscription)                                         // Add a new transcription
		transcriptionProtected.GET("/:transcription_id", a.mediaController.GetTranscriptionByID)                     // Get transcription by ID
		transcriptionProtected.GET("/:transcription_id/user/:userID", a.mediaController.GetTranscriptionByUserID)    // Get transcription by transcription ID and user ID
		transcriptionProtected.GET("/:transcription_id/video/:videoID", a.mediaController.GetTranscriptionByVideoID) // Get transcription by transcription ID and video ID
		transcriptionProtected.GET("/user/:user_id", a.mediaController.ListTranscriptionsByUserID)                   // List transcriptions by user ID
		transcriptionProtected.GET("/video/:video_id", a.mediaController.ListTranscriptionsByVideoID)                // List transcriptions by video ID
		transcriptionProtected.DELETE("/:transcription_id", a.mediaController.DeleteTranscription)                   // Delete transcription by ID
		transcriptionProtected.POST("/generate-upload-url", a.mediaController.GenerateUploadURLForText)              // Generate presigned upload URL
		transcriptionProtected.GET("/:transcription_id/download-url", a.mediaController.GenerateDownloadURLForText)  // Generate presigned download URL
		transcriptionProtected.PUT("/:transcription_id/status", a.mediaController.UpdateTranscriptionStatus)
	}

	audioProtected := r.Group("/audios")
	audioProtected.Use(a.authMiddleware.MustAuth())
	{
		audioProtected.POST("/", a.mediaController.AddAudio)                                  // Add a new audio
		audioProtected.GET("/:audio_id", a.mediaController.GetAudio)                          // Get a specific audio by ID
		audioProtected.DELETE("/:audio_id", a.mediaController.DeleteAudio)                    // Delete an audio
		audioProtected.GET("/user/:user_id", a.mediaController.ListAudiosByUserID)            // Get all audios by user
		audioProtected.GET("/video/:video_id", a.mediaController.ListAudiosByVideoID)         // Get all audios by video
		audioProtected.GET("/:audio_id/user/:user_id", a.mediaController.GetAudioByUser)      // Get specific audio by audio ID and user ID
		audioProtected.GET("/:audio_id/video/:video_id", a.mediaController.GetAudioByVideoID) // Get specific audio by audio ID and video ID
		audioProtected.POST("/generate-presigned-url", a.mediaController.GenerateUploadURL)   // Generate presigned URL for audio upload
		audioProtected.GET("/:audio_id/download-url", a.mediaController.GenerateDownloadURL)  // Generate presigned URL for audio download
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
		public.GET("/:user_id", a.progressController.GetUserProgress)
		public.POST("update-title/:progress_id", a.progressController.UpdateProgressTitle)
		public.POST("delete-progress/:progress_id", a.progressController.DeleteProgress)
	}
}

// RegisterAdminRoutes sets up the routes for all permission-related operations
func (a *AppRouter) RegisterAdminRoutes(r *gin.RouterGroup) {
	protected := r.Group("/admin")
	protected.Use(a.authMiddleware.MustAuth())
	{
		// Config
		protected.GET("/:adminID/config", a.adminController.GetServerConfig)
		protected.POST("/:adminID/config", a.adminController.UpdateServerConfig)

		// Model Options
		protected.GET("/:adminID/models", a.adminController.GetModelList)
		protected.POST("/:adminID/models", a.adminController.AddModelOption)
		protected.PUT("/:adminID/models/:modelOptionID", a.adminController.UpdateModelOption)

		// Monitor
		protected.POST("/:adminID/monitor/media-report", a.adminController.GetMonitorDataType)
		protected.POST("/:adminID/monitor/pipeline-report", a.adminController.GetMonitorPipeline)
		protected.POST("/:adminID/monitor/traffic-report", a.adminController.GetMonitorTraffic)
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

func (a *AppRouter) RegisterPaymentRoutes(r *gin.RouterGroup) {
	// Public endpoints (no auth required)
	public := r.Group("/payment")
	{
		public.GET("/options", a.paymentController.GetPaymentOptions)
	}

	protected := r.Group("/payment")
	protected.Use(a.authMiddleware.MustAuth())
	{
		// User payment endpoints
		protected.POST("/create", a.paymentController.CreatePayment)
		protected.GET("/:payment_id", a.paymentController.GetPayment)
		protected.GET("/transaction/:transaction_id", a.paymentController.GetPaymentByTransactionID)
		protected.GET("/user-payments", a.paymentController.GetUserPayments)
		protected.POST("/confirm/:transaction_id", a.paymentController.ConfirmPayment)
		protected.POST("/cancel/:payment_id", a.paymentController.CancelPayment)

		// Admin endpoints (for monitoring pending payments)
		protected.GET("/pending", a.paymentController.GetPendingPayments)
	}
}

func (a *AppRouter) RegisteVoucherRoutes(r *gin.RouterGroup) {
	protected := r.Group("/voucher")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.POST("/create", a.voucherController.CreateVoucher)
		protected.POST("/use/:code", a.voucherController.UseVoucher)
		protected.PATCH("/:voucherID", a.voucherController.UpdateVoucher)
		protected.GET("/get-all", a.voucherController.GetAllVouchers)
		protected.GET("/:voucherID", a.voucherController.GetVoucherByID)
	}
}

func (a *AppRouter) RegisterTokenRoutes(r *gin.RouterGroup) {
	token := r.Group("/token")
	{
		token.POST("/daily", a.tokenController.ClaimDaily)
		token.POST("/premium", a.tokenController.ClaimPremium)
		token.GET("/claims", a.tokenController.ListClaims)
	}

	premium := r.Group("/premium")
	{
		premium.POST("/add", a.tokenController.AddPremium)
		premium.GET("/list", a.tokenController.ListPremium)
		premium.GET("/:user_id", a.tokenController.CheckPremium)
	}
}

// RegisterSwaggerRoutes sets up the route for Swagger API documentation
func (a *AppRouter) RegisterSwaggerRoutes(r *gin.RouterGroup) {
	// Check if SwaggerRouter is initialized before registering
	if a.swaggerRouter != nil {
		a.swaggerRouter.Register(r)
	}
}

// RegisterMediaRoutes sets up the routes for process-related operations
func (a *AppRouter) RegisterProcessRoutes(r *gin.RouterGroup) {
	processProtected := r.Group("/process")
	processProtected.Use(a.authMiddleware.MustAuth())
	{
		processProtected.POST("/get-all/:user_id", a.processController.GetAllProcess)
	}
}
