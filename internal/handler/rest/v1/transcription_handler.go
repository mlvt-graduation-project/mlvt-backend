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
