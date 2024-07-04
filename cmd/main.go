package main

import (
	"log"
	"mlvt/routes"

	"github.com/gin-gonic/gin"
)

const port = "8080"

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

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
