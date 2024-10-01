package handler

import (
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TranscriptionController struct {
	transcriptionService service.TranscriptionService
}

func NewTranscriptionController(transcriptionService service.TranscriptionService) *TranscriptionController {
	return &TranscriptionController{transcriptionService: transcriptionService}
}

// GenerateUploadURL generates a presigned URL for uploading a transcription file.
// @Summary Generate presigned upload URL
// @Description Generates a presigned URL to upload a transcription file to the storage service.
// @Tags transcriptions
// @Produce json
// @Param file_name query string true "Name of the file to be uploaded"
// @Param file_type query string true "MIME type of the file (e.g., application/json)"
// @Success 200 {object} map[string]string "upload_url"
// @Failure 500 {object} map[string]string "error"
// @Router /transcriptions/generate-upload-url [post]
func (h *TranscriptionController) GenerateUploadURL(c *gin.Context) {
	folder := env.EnvConfig.TranscriptionsFolder
	fileName := c.Query("file_name")
	fileType := c.Query("file_type")

	url, err := h.transcriptionService.GeneratePresignedUploadURL(folder, fileName, fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url})
}

// GenerateDownloadURL generates a presigned URL for downloading a transcription file
// @Summary Generate presigned download URL
// @Description Generates a presigned URL to download a transcription file from the storage service.
// @Tags transcriptions
// @Produce json
// @Param transcription_id path int true "ID of the transcription file"
// @Success 200 {object} map[string]string "download_url"
// @Failure 500 {object} map[string]string "error"
// @Router /transcriptions/{transcription_id}/download-url [get]
func (h *TranscriptionController) GenerateDownloadURL(c *gin.Context) {
	// Parse transcription ID from the URL path
	transcriptionID, err := strconv.ParseUint(c.Param("transcription_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transcription ID"})
		return
	}

	// Call the service to generate the presigned download URL
	downloadURL, err := h.transcriptionService.GeneratePresignedDownloadURL(transcriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the presigned download URL
	c.JSON(http.StatusOK, gin.H{"download_url": downloadURL})
}

// AddTranscription adds a new transcription record.
// @Summary Add transcription
// @Description Adds a new transcription file's metadata to the system.
// @Tags transcriptions
// @Accept json
// @Produce json
// @Param transcription body entity.Transcription true "Transcription object"
// @Success 201 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /transcriptions [post]
func (h *TranscriptionController) AddTranscription(c *gin.Context) {
	var transcription entity.Transcription
	if err := c.ShouldBindJSON(&transcription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.transcriptionService.CreateTranscription(&transcription); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transcription added successfully"})
}

// GetTranscriptionByID retrieves a transcription by ID and generates a presigned download URL.
// @Summary Get transcription by ID
// @Description Retrieves a transcription and generates a presigned download URL for it.
// @Tags transcriptions
// @Produce json
// @Param transcription_id path int true "ID of the transcription file"
// @Success 200 {object} map[string]interface{} "transcription, download_url"
// @Failure 404 {object} map[string]string "error"
// @Router /transcriptions/{transcription_id} [get]
func (h *TranscriptionController) GetTranscriptionByID(c *gin.Context) {
	transcriptionID, _ := strconv.ParseUint(c.Param("transcription_id"), 10, 64)

	transcription, downloadURL, err := h.transcriptionService.GetTranscriptionByID(transcriptionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transcription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transcription": transcription,
		"download_url":  downloadURL,
	})
}

// GetTranscriptionByUserID retrieves a transcription by transcription ID and User ID, generating a download URL.
// @Summary Get transcription by User ID and transcription ID
// @Description Retrieves a transcription for a specific user and generates a presigned download URL.
// @Tags transcriptions
// @Produce json
// @Param transcriptionID path int true "ID of the transcription file"
// @Param userID path int true "ID of the user"
// @Success 200 {object} map[string]interface{} "transcription, download_url"
// @Failure 500 {object} map[string]string "error"
// @Router /transcriptions/{transcriptionID}/user/{userID} [get]
func (h *TranscriptionController) GetTranscriptionByUserID(c *gin.Context) {
	transcriptionID, _ := strconv.ParseUint(c.Param("transcriptionID"), 10, 64)
	userID, _ := strconv.ParseUint(c.Param("userID"), 10, 64)

	transcription, downloadURL, err := h.transcriptionService.GetTranscriptionByIDAndUserID(transcriptionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transcription": transcription,
		"download_url":  downloadURL,
	})
}

// GetTranscriptionByVideoID retrieves a transcription by transcription ID and Video ID, generating a download URL.
// @Summary Get transcription by Video ID and transcription ID
// @Description Retrieves a transcription for a specific video and generates a presigned download URL.
// @Tags transcriptions
// @Produce json
// @Param transcriptionID path int true "ID of the transcription file"
// @Param videoID path int true "ID of the video"
// @Success 200 {object} map[string]interface{} "transcription, download_url"
// @Failure 500 {object} map[string]string "error"
// @Router /transcriptions/{transcriptionID}/video/{videoID} [get]
func (h *TranscriptionController) GetTranscriptionByVideoID(c *gin.Context) {
	transcriptionID, _ := strconv.ParseUint(c.Param("transcriptionID"), 10, 64)
	videoID, _ := strconv.ParseUint(c.Param("videoID"), 10, 64)

	transcription, downloadURL, err := h.transcriptionService.GetTranscriptionByIDAndVideoID(transcriptionID, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transcription": transcription,
		"download_url":  downloadURL,
	})
}

// ListTranscriptionsByUserID lists all transcriptions for a specific user.
// @Summary List transcriptions by User ID
// @Description Retrieves all transcriptions belonging to a specific user.
// @Tags transcriptions
// @Produce json
// @Param user_id path int true "ID of the user"
// @Success 200 {array} entity.Transcription "transcriptions"
// @Failure 500 {object} map[string]string "error"
// @Router /transcriptions/user/{user_id} [get]
func (h *TranscriptionController) ListTranscriptionsByUserID(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	transcriptions, err := h.transcriptionService.ListTranscriptionsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transcriptions": transcriptions})
}

// ListTranscriptionsByVideoID lists all transcriptions for a specific video.
// @Summary List transcriptions by Video ID
// @Description Retrieves all transcriptions belonging to a specific video.
// @Tags transcriptions
// @Produce json
// @Param video_id path int true "ID of the video"
// @Success 200 {array} entity.Transcription "transcriptions"
// @Failure 500 {object} map[string]string "error"
// @Router /transcriptions/video/{video_id} [get]
func (h *TranscriptionController) ListTranscriptionsByVideoID(c *gin.Context) {
	videoID, _ := strconv.ParseUint(c.Param("video_id"), 10, 64)

	transcriptions, err := h.transcriptionService.ListTranscriptionsByVideoID(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transcriptions": transcriptions})
}

// DeleteTranscription deletes a transcription by its ID.
// @Summary Delete transcription by ID
// @Description Deletes a transcription record from the system.
// @Tags transcriptions
// @Param transcription_id path int true "ID of the transcription file"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /transcriptions/{transcription_id} [delete]
func (h *TranscriptionController) DeleteTranscription(c *gin.Context) {
	transcriptionID, _ := strconv.ParseUint(c.Param("transcription_id"), 10, 64)

	if err := h.transcriptionService.DeleteTranscription(transcriptionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transcription deleted successfully"})
}
