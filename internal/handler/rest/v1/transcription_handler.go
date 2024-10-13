package handler

import (
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/pkg/response"
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

// GenerateUploadURL godoc
// @Summary Generate presigned upload URL
// @Description Generates a presigned URL to upload a transcription file to the storage service.
// @Tags transcriptions
// @Produce json
// @Param file_name query string true "Name of the file to be uploaded"
// @Param file_type query string true "MIME type of the file (e.g., application/json)"
// @Success 200 {object} response.UploadURLResponse "upload_url"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /transcriptions/generate-upload-url [post]
func (h *TranscriptionController) GenerateUploadURL(c *gin.Context) {
	folder := env.EnvConfig.TranscriptionsFolder
	fileName := c.Query("file_name")
	fileType := c.Query("file_type")

	url, err := h.transcriptionService.GeneratePresignedUploadURL(folder, fileName, fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.UploadURLResponse{UploadURL: url})
}

// GenerateDownloadURL godoc
// @Summary Generate presigned download URL
// @Description Generates a presigned URL to download a transcription file from the storage service.
// @Tags transcriptions
// @Produce json
// @Param transcription_id path uint64 true "ID of the transcription file"
// @Success 200 {object} response.DownloadURLResponse "download_url"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /transcriptions/{transcription_id}/download-url [get]
func (h *TranscriptionController) GenerateDownloadURL(c *gin.Context) {
	// Parse transcription ID from the URL path
	transcriptionIDStr := c.Param("transcription_id")
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transcription ID"})
		return
	}

	// Call the service to generate the presigned download URL
	downloadURL, err := h.transcriptionService.GeneratePresignedDownloadURL(transcriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Return the presigned download URL
	c.JSON(http.StatusOK, response.DownloadURLResponse{DownloadURL: downloadURL})
}

// AddTranscription godoc
// @Summary Add transcription
// @Description Adds a new transcription file's metadata to the system.
// @Tags transcriptions
// @Accept json
// @Produce json
// @Param transcription body entity.Transcription true "Transcription object"
// @Success 201 {object} response.MessageResponse "message"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /transcriptions [post]
func (h *TranscriptionController) AddTranscription(c *gin.Context) {
	var transcription entity.Transcription
	if err := c.ShouldBindJSON(&transcription); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.transcriptionService.CreateTranscription(&transcription); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.MessageResponse{Message: "Transcription added successfully"})
}

// GetTranscriptionByID godoc
// @Summary Get transcription by ID
// @Description Retrieves a transcription and generates a presigned download URL for it.
// @Tags transcriptions
// @Produce json
// @Param transcription_id path uint64 true "ID of the transcription file"
// @Success 200 {object} response.TranscriptionResponse "transcription, download_url"
// @Failure 404 {object} response.ErrorResponse "error"
// @Router /transcriptions/{transcription_id} [get]
func (h *TranscriptionController) GetTranscriptionByID(c *gin.Context) {
	transcriptionIDStr := c.Param("transcription_id")
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transcription ID"})
		return
	}

	transcription, downloadURL, err := h.transcriptionService.GetTranscriptionByID(transcriptionID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "transcription not found"})
		return
	}

	c.JSON(http.StatusOK, response.TranscriptionResponse{
		Transcription: *transcription,
		DownloadURL:   downloadURL,
	})
}

// GetTranscriptionByUserID godoc
// @Summary Get transcription by User ID and transcription ID
// @Description Retrieves a transcription for a specific user and generates a presigned download URL.
// @Tags transcriptions
// @Produce json
// @Param transcriptionID path uint64 true "ID of the transcription file"
// @Param userID path uint64 true "ID of the user"
// @Success 200 {object} response.TranscriptionResponse "transcription, download_url"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /transcriptions/{transcriptionID}/user/{userID} [get]
func (h *TranscriptionController) GetTranscriptionByUserID(c *gin.Context) {
	transcriptionIDStr := c.Param("transcriptionID")
	userIDStr := c.Param("userID")

	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transcription ID"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user ID"})
		return
	}

	transcription, downloadURL, err := h.transcriptionService.GetTranscriptionByIDAndUserID(transcriptionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.TranscriptionResponse{
		Transcription: *transcription,
		DownloadURL:   downloadURL,
	})
}

// GetTranscriptionByVideoID godoc
// @Summary Get transcription by Video ID and transcription ID
// @Description Retrieves a transcription for a specific video and generates a presigned download URL.
// @Tags transcriptions
// @Produce json
// @Param transcriptionID path uint64 true "ID of the transcription file"
// @Param videoID path uint64 true "ID of the video"
// @Success 200 {object} response.TranscriptionResponse "transcription, download_url"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /transcriptions/{transcriptionID}/video/{videoID} [get]
func (h *TranscriptionController) GetTranscriptionByVideoID(c *gin.Context) {
	transcriptionIDStr := c.Param("transcriptionID")
	videoIDStr := c.Param("videoID")

	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transcription ID"})
		return
	}

	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid video ID"})
		return
	}

	transcription, downloadURL, err := h.transcriptionService.GetTranscriptionByIDAndVideoID(transcriptionID, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.TranscriptionResponse{
		Transcription: *transcription,
		DownloadURL:   downloadURL,
	})
}

// ListTranscriptionsByUserID godoc
// @Summary List transcriptions by User ID
// @Description Retrieves all transcriptions belonging to a specific user.
// @Tags transcriptions
// @Produce json
// @Param user_id path uint64 true "ID of the user"
// @Success 200 {object} response.TranscriptionsResponse "transcriptions"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /transcriptions/user/{user_id} [get]
func (h *TranscriptionController) ListTranscriptionsByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user ID"})
		return
	}

	transcriptions, err := h.transcriptionService.ListTranscriptionsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response.TranscriptionsResponse{Transcriptions: transcriptions})
}

// ListTranscriptionsByVideoID godoc
// @Summary List transcriptions by Video ID
// @Description Retrieves all transcriptions belonging to a specific video.
// @Tags transcriptions
// @Produce json
// @Param video_id path uint64 true "ID of the video"
// @Success 200 {object} response.TranscriptionsResponse "transcriptions"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /transcriptions/video/{video_id} [get]
func (h *TranscriptionController) ListTranscriptionsByVideoID(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid video ID"})
		return
	}

	transcriptions, err := h.transcriptionService.ListTranscriptionsByVideoID(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response.TranscriptionsResponse{Transcriptions: transcriptions})
}

// DeleteTranscription godoc
// @Summary Delete transcription by ID
// @Description Deletes a transcription record from the system.
// @Tags transcriptions
// @Param transcription_id path uint64 true "ID of the transcription file"
// @Success 200 {object} response.MessageResponse "message"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /transcriptions/{transcription_id} [delete]
func (h *TranscriptionController) DeleteTranscription(c *gin.Context) {
	transcriptionIDStr := c.Param("transcription_id")
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transcription ID"})
		return
	}

	if err := h.transcriptionService.DeleteTranscription(transcriptionID); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Transcription deleted successfully"})
}
