package handler

import (
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AudioController struct {
	audioService service.AudioService
}

func NewAudioController(audioService service.AudioService) *AudioController {
	return &AudioController{
		audioService: audioService,
	}
}

// GenerateUploadURL generates a presigned URL for uploading an audio file.
// @Summary Generate presigned upload URL
// @Description Generates a presigned URL to upload an audio file to the storage service.
// @Tags audios
// @Produce  json
// @Param file_name query string true "Name of the file to be uploaded"
// @Param file_type query string true "MIME type of the file (e.g., audio/mpeg)"
// @Success 200 {object} map[string]string "upload_url"
// @Failure 500 {object} map[string]string "error"
// @Router /audios/generate-upload-url [get]
func (h *AudioController) GenerateUploadURL(c *gin.Context) {
	folder := env.EnvConfig.AudioFolder
	fileName := c.Query("file_name")
	fileType := c.Query("file_type")

	url, err := h.audioService.GeneratePresignedUploadURL(folder, fileName, fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url})
}

// GenerateDownloadURL generates a presigned URL for downloading an audio file.
// @Summary Generate presigned download URL
// @Description Generates a presigned URL to download an audio file from the storage service.
// @Tags audios
// @Produce  json
// @Param audioID path int true "ID of the audio file"
// @Success 200 {object} map[string]string "download_url"
// @Failure 500 {object} map[string]string "error"
// @Router /audios/{audioID}/download-url [get]
func (h *AudioController) GenerateDownloadURL(c *gin.Context) {
	audioID, _ := strconv.ParseUint(c.Param("audioID"), 10, 64)

	url, err := h.audioService.GeneratePresignedDownloadURL(audioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"download_url": url})
}

// AddAudio adds a new audio record.
// @Summary Add audio
// @Description Adds a new audio file's metadata to the system.
// @Tags audios
// @Accept  json
// @Produce  json
// @Param audio body entity.Audio true "Audio object"
// @Success 201 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /audios [post]
func (h *AudioController) AddAudio(c *gin.Context) {
	var audio entity.Audio
	if err := c.ShouldBindJSON(&audio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.audioService.CreateAudio(&audio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Audio added successfully"})
}

// GetAudio retrieves an audio file and its presigned download URL by audio ID.
// @Summary Get audio by ID
// @Description Retrieves an audio file's metadata and generates a presigned download URL.
// @Tags audios
// @Produce  json
// @Param audio_id path int true "ID of the audio file"
// @Success 200 {object} map[string]interface{} "audio, download_url"
// @Failure 404 {object} map[string]string "error"
// @Router /audios/{audio_id} [get]
func (h *AudioController) GetAudio(c *gin.Context) {
	audioID, _ := strconv.ParseUint(c.Param("audio_id"), 10, 64)

	audio, url, err := h.audioService.GetAudioByID(audioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audio not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"audio": audio, "download_url": url})
}

// GetAudioByUser retrieves an audio file by audio ID and user ID, and generates a presigned download URL.
// @Summary Get audio by user and audio ID
// @Description Retrieves an audio file for a specific user and generates a presigned download URL.
// @Tags audios
// @Produce  json
// @Param audioID path int true "ID of the audio file"
// @Param userID path int true "ID of the user"
// @Success 200 {object} map[string]interface{} "audio, download_url"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /audios/{audioID}/user/{userID} [get]
func (h *AudioController) GetAudioByUser(c *gin.Context) {
	// Parse audio ID from the URL path
	audioID, err := strconv.ParseUint(c.Param("audioID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio ID"})
		return
	}

	// Parse user ID from the URL path
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the service to get the audio and the presigned download URL
	audio, downloadURL, err := h.audioService.GetAudioByIDAndUserID(audioID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the audio metadata and the presigned URL
	c.JSON(http.StatusOK, gin.H{
		"audio":        audio,
		"download_url": downloadURL,
	})
}

// ListAudiosByUserID lists all audio files for a specific user.
// @Summary List audios by user ID
// @Description Retrieves all audio files belonging to a specific user.
// @Tags audios
// @Produce  json
// @Param user_id path int true "ID of the user"
// @Success 200 {array} entity.Audio "audios"
// @Failure 500 {object} map[string]string "error"
// @Router /audios/user/{user_id} [get]
func (h *AudioController) ListAudiosByUserID(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	audios, err := h.audioService.ListAudiosByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"audios": audios})
}

// GetAudioByVideo retrieves an audio file by audio ID and video ID, and generates a presigned download URL.
// @Summary Get audio by video and audio ID
// @Description Retrieves an audio file for a specific video and generates a presigned download URL.
// @Tags audios
// @Produce  json
// @Param audioID path int true "ID of the audio file"
// @Param videoID path int true "ID of the video"
// @Success 200 {object} map[string]interface{} "audio, download_url"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /audios/{audioID}/video/{videoID} [get]
func (h *AudioController) GetAudioByVideo(c *gin.Context) {
	// Parse audio ID from the URL path
	audioID, err := strconv.ParseUint(c.Param("audioID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio ID"})
		return
	}

	// Parse video ID from the URL path
	videoID, err := strconv.ParseUint(c.Param("videoID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Call the service to get the audio and the presigned download URL
	audio, downloadURL, err := h.audioService.GetAudioByVideoID(videoID, audioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the audio metadata and the presigned URL
	c.JSON(http.StatusOK, gin.H{
		"audio":        audio,
		"download_url": downloadURL,
	})
}

// ListAudiosByVideoID lists all audio files associated with a specific video.
// @Summary List audios by video ID
// @Description Retrieves all audio files belonging to a specific video.
// @Tags audios
// @Produce  json
// @Param video_id path int true "ID of the video"
// @Success 200 {array} entity.Audio "audios"
// @Failure 500 {object} map[string]string "error"
// @Router /audios/video/{video_id} [get]
func (h *AudioController) ListAudiosByVideoID(c *gin.Context) {
	videoID, _ := strconv.ParseUint(c.Param("video_id"), 10, 64)

	audios, err := h.audioService.ListAudiosByVideoID(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"audios": audios})
}

// DeleteAudio deletes an audio file by its ID.
// @Summary Delete audio by ID
// @Description Deletes an audio file from the system.
// @Tags audios
// @Param id path int true "ID of the audio file"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /audios/{id} [delete]
func (h *AudioController) DeleteAudio(c *gin.Context) {
	// Parse audio ID from the URL path
	audioID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio ID"})
		return
	}

	// Call the service to delete the audio
	if err := h.audioService.DeleteAudio(audioID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete audio"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Audio deleted successfully"})
}
