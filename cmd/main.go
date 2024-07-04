package main

import (
	"log"
	"mlvt/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const port = "8080"

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.MaxMultipartMemory = 32 << 20

	//Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET, PUT, POST, DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	trustedProxies := []string{"192.168.0.1", "10.0.0.1"}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	routes.SetUpRoutes(r)

	log.Println("Connecting to: ", port)

	r.Run(":" + port)
}
