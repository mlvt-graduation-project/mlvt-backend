package routes

import (
	"mlvt/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	router.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.POST("/full_pipeline/upload", handlers.UploadVideos)
}
