package media_handler

import (
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GenerateUploadURLForText godoc
// @Summary Generate presigned upload URL
// @Description Generates a presigned URL to upload a transcription file to the storage service.
// @Tags transcriptions
// @Produce json
// @Param file_name query string true "Name of the file to be uploaded"
// @Param file_type query string true "MIME type of the file (e.g., application/json)"
// @Success 200 {object} response.UploadURLResponse "upload_url"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /transcriptions/generate-upload-url [post]
func (h *MediaController) GenerateUploadURLForText(c *gin.Context) {
	folder := env.EnvConfig.TranscriptionsFolder
	fileName := c.Query("file_name")
	fileType := c.Query("file_type")

	url, err := h.mediaService.GeneratePresignedUploadURLForText(folder, fileName, fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.UploadURLResponse{UploadURL: url})
}

// GenerateDownloadURLForText godoc
// @Summary Generate presigned download URL
// @Description Generates a presigned URL to download a transcription file from the storage service.
// @Tags transcriptions
// @Produce json
// @Param transcription_id path uint64 true "ID of the transcription file"
// @Success 200 {object} response.DownloadURLResponse "download_url"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /transcriptions/{transcription_id}/download-url [get]
func (h *MediaController) GenerateDownloadURLForText(c *gin.Context) {
	// Parse transcription ID from the URL path
	transcriptionIDStr := c.Param("transcription_id")
	log.Warnf("extract from param: %s", transcriptionIDStr)
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transcription ID"})
		return
	}

	// Call the service to generate the presigned download URL
	downloadURL, err := h.mediaService.GeneratePresignedDownloadURLForText(transcriptionID)
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
func (h *MediaController) AddTranscription(c *gin.Context) {
	var transcription entity.Transcription
	if err := c.ShouldBindJSON(&transcription); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := h.mediaService.CreateTranscription(&transcription, false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.MessageCreateResponseWithID{
		Message: "Video added successfully",
		Id:      id,
	})
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
func (h *MediaController) GetTranscriptionByID(c *gin.Context) {
	transcriptionIDStr := c.Param("transcription_id")
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transcription ID"})
		return
	}

	transcription, downloadURL, err := h.mediaService.GetTranscriptionByID(transcriptionID)
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
func (h *MediaController) GetTranscriptionByUserID(c *gin.Context) {
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

	transcription, downloadURL, err := h.mediaService.GetTranscriptionByIDAndUserID(transcriptionID, userID)
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
func (h *MediaController) GetTranscriptionByVideoID(c *gin.Context) {
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

	transcription, downloadURL, err := h.mediaService.GetTranscriptionByIDAndVideoID(transcriptionID, videoID)
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
func (h *MediaController) ListTranscriptionsByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user ID"})
		return
	}

	transcriptions, err := h.mediaService.ListTranscriptionsByUserID(userID)
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
func (h *MediaController) ListTranscriptionsByVideoID(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid video ID"})
		return
	}

	transcriptions, err := h.mediaService.ListTranscriptionsByVideoID(videoID)
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
func (h *MediaController) DeleteTranscription(c *gin.Context) {
	transcriptionIDStr := c.Param("transcription_id")
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transcription ID"})
		return
	}

	if err := h.mediaService.DeleteTranscription(transcriptionID); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Transcription deleted successfully"})
}

type UpdateTranscriptionStatusRequest struct {
	Status entity.StatusEntity `json:"status"`
}

// UpdateTranscriptionStatus godoc
// @Summary Update the status of a transcription
// @Description Update the status of a specific Transcription by its ID
// @Tags transcriptions
// @Accept  json
// @Produce  json
// @Param   transcription_id path     uint64 true "Transcription ID"
// @Param   status   body     UpdateTranscriptionStatusRequest true "New status"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transcriptions/{transcription_id}/status [put]
func (h *MediaController) UpdateTranscriptionStatus(c *gin.Context) {
	transcriptionIDStr := c.Param("transcription_id")
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid Transcription ID"})
		return
	}

	var req UpdateTranscriptionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid input"})
		return
	}

	err = h.mediaService.UpdateTranscriptionStatus(transcriptionID, req.Status)
	if err != nil {
		if err.Error() == "no transcription found with id "+strconv.FormatUint(transcriptionID, 10) {
			c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "video not found"})
		} else {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "status updated successfully"})
}
