package initialize

import (
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/server/http"
	"mlvt/internal/router"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitServer configures and returns the HTTP server instance.
func InitServer(appRouter *router.AppRouter) *http.Server {
	// Create a new Gin router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins during development
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Allow credentials like cookies
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	api := r.Group("/api")
	appRouter.RegisterUserRoutes(api)
	appRouter.RegisterProcessRoutes(api)
	appRouter.RegisterMediaRoutes(api)
	appRouter.RegiserMlvtRoutes(api)
	appRouter.RegisterProgressRoutes(api)
	appRouter.RegisterPingStatusRoutes(api)
	appRouter.RegisterAdminRoutes(api)
	appRouter.RegisterWalletRoutes(api)
	appRouter.RegisterPaymentRoutes(api)
	appRouter.RegisteVoucherRoutes(api)
	appRouter.RegisterTokenRoutes(api)
	appRouter.RegisterSwaggerRoutes(r.Group("/"))

	// Create the HTTP server
	addr := ":" + env.EnvConfig.ServerPort
	server := http.NewServer(r, addr)

	return server
}
