package handler

import (
	"mlvt/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TranscriptionController struct {
	Service service.TranscriptionService
}

func NewTranscriptionController(service service.TranscriptionService) *TranscriptionController {
	return &TranscriptionController{Service: service}
}

// GeneratePresignedURL godoc
// @Summary Generate a presigned URL for a video
// @Description Generate a presigned URL to upload a transcription for a specific video.
// @Tags transcription
// @Accept json
// @Produce json
// @Param videoID path uint64 true "Video ID"
// @Success 200 {object} map[string]string "presigned URL"
// @Failure 400 {object} map[string]string "Invalid video ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transcriptions/presigned-url/{videoID} [get]
func (h *TranscriptionController) GeneratePresignedURL(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("videoID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	url, err := h.Service.GeneratePresignedURLForTranscription(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// CreateTranscription godoc
// @Summary Create a new transcription
// @Description Create a new transcription with provided details.
// @Tags transcription
// @Accept json
// @Produce json
// @Param transcription body struct { VideoID uint64 "video_id"; UserID uint64 "user_id"; Lang string; TranscriptionText string } true "Transcription Details"
// @Success 200 {object} entity.Transcription
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transcriptions [post]
func (h *TranscriptionController) CreateTranscription(c *gin.Context) {
	var req struct {
		VideoID           uint64 `json:"video_id"`
		UserID            uint64 `json:"user_id"`
		Lang              string `json:"lang"`
		TranscriptionText string `json:"transcription_text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	tx, err := h.Service.CreateTranscription(req.VideoID, req.UserID, req.Lang, req.TranscriptionText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

// GetTranscription godoc
// @Summary Get a single transcription
// @Description Get a single transcription by user ID and transcription ID.
// @Tags transcription
// @Accept json
// @Produce json
// @Param userID path uint64 true "User ID"
// @Param transcriptionID path uint64 true "Transcription ID"
// @Success 200 {object} entity.Transcription
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transcriptions/{userID}/{transcriptionID} [get]
func (h *TranscriptionController) GetTranscription(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	transcriptionID, err := strconv.ParseUint(c.Param("transcriptionID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	tx, err := h.Service.GetTranscription(userID, transcriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

// ListTranscriptionsByUser godoc
// @Summary List transcriptions by user
// @Description List all transcriptions associated with a user ID.
// @Tags transcription
// @Accept json
// @Produce json
// @Param userID path uint64 true "User ID"
// @Success 200 {array} entity.Transcription
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transcriptions/by-user/{userID} [get]
func (h *TranscriptionController) ListTranscriptionsByUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	transcriptions, err := h.Service.ListTranscriptionsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transcriptions)
}

// GetTranscriptionByVideo godoc
// @Summary Get a transcription by video
// @Description Get a specific transcription by video ID and transcription ID.
// @Tags transcription
// @Accept json
// @Produce json
// @Param videoID path uint64 true "Video ID"
// @Param transcriptionID path uint64 true "Transcription ID"
// @Success 200 {object} entity.Transcription
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transcriptions/by-video/{videoID}/{transcriptionID} [get]
func (h *TranscriptionController) GetTranscriptionByVideo(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("videoID"), 10, 64)
	transcriptionID, err := strconv.ParseUint(c.Param("transcriptionID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	tx, err := h.Service.GetTranscriptionByVideo(videoID, transcriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

// ListTranscriptionsByVideo godoc
// @Summary List transcriptions by video
// @Description List all transcriptions associated with a video ID.
// @Tags transcription
// @Accept json
// @Produce json
// @Param videoID path uint64 true "Video ID"
// @Success 200 {array} entity.Transcription
// @Failure 400 {object} map[string]string "Invalid video ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transcriptions/by-video/{videoID} [get]
func (h *TranscriptionController) ListTranscriptionsByVideo(c *gin.Context) {
	videoID, err := strconv.ParseUint(c.Param("videoID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	transcriptions, err := h.Service.ListTranscriptionsByVideo(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transcriptions)
}

// DeleteTranscription godoc
// @Summary Delete a transcription
// @Description Delete a transcription by transcription ID.
// @Tags transcription
// @Accept json
// @Produce json
// @Param transcriptionID path uint64 true "Transcription ID"
// @Success 200
// @Failure 400 {object} map[string]string "Invalid transcription ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transcriptions/{transcriptionID} [delete]
func (h *TranscriptionController) DeleteTranscription(c *gin.Context) {
	transcriptionID, err := strconv.ParseUint(c.Param("transcriptionID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transcription ID"})
		return
	}

	err = h.Service.DeleteTranscription(transcriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
