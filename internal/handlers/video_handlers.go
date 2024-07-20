package handlers

import (
	"errors"
	"log"
	"mlvt/internal/service"
	utils "mlvt/pkg"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService *service.VideoService
}

func NewVideoHandler(videoService *service.VideoService) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
	}
}

func (h *VideoHandler) UploadVideos(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusBadRequest)
		return
	}

	userID := c.Query("user_id")
	if userID == "" {
		utils.ErrorJSON(c, errors.New("user ID is required"), http.StatusBadRequest)
		return
	}

	tempPath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		utils.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempPath)

	video, err := h.videoService.ProcessVideo(tempPath)
	if err != nil {
		log.Printf("error processing video: %v", err)
		utils.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	video.Type = file.Header.Get("Content-Type")
	video.UserID = userID

	bucketName := ""
	if err := h.videoService.UploadAndSaveVideo(video, bucketName); err != nil {
		log.Printf("error uploading and saving video: %v", err)
		utils.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	log.Printf("received video - ID: %s, Duration: %d sec, Type: %s, UploadedAt: %s, Size: %d bytes, UserID: %s, Type: %s", video.ID, video.Duration, video.FilePath, video.UploadedAt.Format(time.RFC3339), video.Size, video.UserID, video.Type)

	response := utils.JSONRespone{
		Error:   false,
		Message: "Video uploaded successfully",
		Data: gin.H{
			"id":         video.ID,
			"duration":   video.Duration,
			"uploadedAt": video.UploadedAt.Format(time.RFC3339),
			"size":       video.Size,
			"type":       video.Type,
		},
	}

	utils.WriteJSON(c, http.StatusOK, response)
}

func (h *VideoHandler) GetUserVideos(c *gin.Context) {
	userID := c.Param("userID")

	videos, err := h.videoService.GetUserVideos(userID)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	response := utils.JSONRespone{
		Error:   false,
		Message: "Videos retrieved successfully",
		Data:    videos,
	}

	utils.WriteJSON(c, http.StatusOK, response)
}
