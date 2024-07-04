package handlers

import (
	"log"
	internal "mlvt/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadVideos(c *gin.Context) {
	file, err := c.FormFile("video")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no video file provided"})
	}

	video := internal.Video{
		FilePath:   file.Filename, // MIME type as video type
		UploadedAt: time.Now(),
		Size:       file.Size,
		//UserId
	}

	log.Printf("Received video - ID: %d, Duration: %d sec, Type: %s, UploadedAt: %s, Size: %d bytes, UserID: %s",
		video.ID, video.Duration, video.FilePath, video.UploadedAt.Format(time.RFC3339), video.Size, video.UserID)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "video uploaded successfully"})

}
