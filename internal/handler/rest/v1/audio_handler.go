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

type AudioController struct {
	audioService service.AudioService
}

func NewAudioController(audioService service.AudioService) *AudioController {
	return &AudioController{
		audioService: audioService,
	}
}

// GenerateUploadURL godoc
// @Summary Generate presigned upload URL
// @Description Generates a presigned URL to upload an audio file to the storage service.
// @Tags audios
// @Produce json
// @Param file_name query string true "Name of the file to be uploaded"
// @Param file_type query string true "MIME type of the file (e.g., audio/mpeg)"
// @Success 200 {object} response.UploadURLResponse "upload_url"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /audios/generate-upload-url [get]
func (h *AudioController) GenerateUploadURL(c *gin.Context) {
	folder := env.EnvConfig.AudioFolder
	fileName := c.Query("file_name")
	fileType := c.Query("file_type")

	// Validate input parameters
	if fileName == "" || fileType == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "file_name and file_type are required"})
		return
	}

	url, err := h.audioService.GeneratePresignedUploadURL(folder, fileName, fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.UploadURLResponse{UploadURL: url})
}

// GenerateDownloadURL godoc
// @Summary Generate presigned download URL
// @Description Generates a presigned URL to download an audio file from the storage service.
// @Tags audios
// @Produce json
// @Param audio_id path uint64 true "ID of the audio file"
// @Success 200 {object} response.DownloadURLResponse "download_url"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /audios/{audio_id}/download-url [get]
func (h *AudioController) GenerateDownloadURL(c *gin.Context) {
	// Parse audio ID from the URL path
	audioIDStr := c.Param("audio_id")
	audioID, err := strconv.ParseUint(audioIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid audio ID"})
		return
	}

	// Call the service to generate the presigned download URL
	downloadURL, err := h.audioService.GeneratePresignedDownloadURL(audioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Return the presigned download URL
	c.JSON(http.StatusOK, response.DownloadURLResponse{DownloadURL: downloadURL})
}

// AddAudio godoc
// @Summary Add audio
// @Description Adds a new audio file's metadata to the system.
// @Tags audios
// @Accept json
// @Produce json
// @Param audio body entity.Audio true "Audio object"
// @Success 201 {object} response.MessageResponse "message"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /audios [post]
func (h *AudioController) AddAudio(c *gin.Context) {
	var audio entity.Audio
	if err := c.ShouldBindJSON(&audio); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.audioService.CreateAudio(&audio); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.MessageResponse{Message: "Audio added successfully"})
}

// GetAudio godoc
// @Summary Get audio by ID
// @Description Retrieves an audio file's metadata and generates a presigned download URL.
// @Tags audios
// @Produce json
// @Param audio_id path uint64 true "ID of the audio file"
// @Success 200 {object} response.AudioResponse "audio, download_url"
// @Failure 404 {object} response.ErrorResponse "error"
// @Router /audios/{audio_id} [get]
func (h *AudioController) GetAudio(c *gin.Context) {
	audioIDStr := c.Param("audio_id")
	audioID, err := strconv.ParseUint(audioIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid audio ID"})
		return
	}

	audio, downloadURL, err := h.audioService.GetAudioByID(audioID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Audio not found"})
		return
	}

	c.JSON(http.StatusOK, response.AudioResponse{
		Audio:       *audio,
		DownloadURL: downloadURL,
	})
}

// GetAudioByUserID godoc
// @Summary Get audio by user and audio ID
// @Description Retrieves an audio file for a specific user and generates a presigned download URL.
// @Tags audios
// @Produce json
// @Param audioID path uint64 true "ID of the audio file"
// @Param userID path uint64 true "ID of the user"
// @Success 200 {object} response.AudioResponse "audio, download_url"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /audios/{audioID}/user/{userID} [get]
func (h *AudioController) GetAudioByUser(c *gin.Context) {
	// Parse audio ID from the URL path
	audioIDStr := c.Param("audioID")
	audioID, err := strconv.ParseUint(audioIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid audio ID"})
		return
	}

	// Parse user ID from the URL path
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Call the service to get the audio and the presigned download URL
	audio, downloadURL, err := h.audioService.GetAudioByIDAndUserID(audioID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Return the audio metadata and the presigned URL
	c.JSON(http.StatusOK, response.AudioResponse{
		Audio:       *audio,
		DownloadURL: downloadURL,
	})
}

// ListAudiosByUserID godoc
// @Summary List audios by user ID
// @Description Retrieves all audio files belonging to a specific user.
// @Tags audios
// @Produce json
// @Param user_id path uint64 true "ID of the user"
// @Success 200 {object} response.AudiosResponse "audios"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /audios/user/{user_id} [get]
func (h *AudioController) ListAudiosByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user ID"})
		return
	}

	transcriptions, err := h.audioService.ListAudiosByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response.AudiosResponse{Audios: transcriptions})
}

// GetAudioByVideoID godoc
// @Summary Get audio by video and audio ID
// @Description Retrieves an audio file for a specific video and generates a presigned download URL.
// @Tags audios
// @Produce json
// @Param audioID path uint64 true "ID of the audio file"
// @Param videoID path uint64 true "ID of the video"
// @Success 200 {object} response.AudioResponse "audio, download_url"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /audios/{audioID}/video/{videoID} [get]
func (h *AudioController) GetAudioByVideoID(c *gin.Context) {
	// Parse audio ID from the URL path
	audioIDStr := c.Param("audioID")
	audioID, err := strconv.ParseUint(audioIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid audio ID"})
		return
	}

	// Parse video ID from the URL path
	videoIDStr := c.Param("videoID")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid video ID"})
		return
	}

	// Call the service to get the audio and the presigned download URL
	audio, downloadURL, err := h.audioService.GetAudioByVideoID(videoID, audioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Return the audio metadata and the presigned URL
	c.JSON(http.StatusOK, response.AudioResponse{
		Audio:       *audio,
		DownloadURL: downloadURL,
	})
}

// ListAudiosByVideoID godoc
// @Summary List audios by Video ID
// @Description Retrieves all audio files belonging to a specific video.
// @Tags audios
// @Produce json
// @Param video_id path uint64 true "ID of the video"
// @Success 200 {object} response.AudiosResponse "audios"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /audios/video/{video_id} [get]
func (h *AudioController) ListAudiosByVideoID(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid video ID"})
		return
	}

	audios, err := h.audioService.ListAudiosByVideoID(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response.AudiosResponse{Audios: audios})
}

// DeleteAudio godoc
// @Summary Delete audio by ID
// @Description Deletes an audio file from the system.
// @Tags audios
// @Param audio_id path uint64 true "ID of the audio file"
// @Success 200 {object} response.MessageResponse "message"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /audios/{audio_id} [delete]
func (h *AudioController) DeleteAudio(c *gin.Context) {
	// Parse audio ID from the URL path
	audioIDStr := c.Param("audio_id")
	audioID, err := strconv.ParseUint(audioIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid audio ID"})
		return
	}

	// Call the service to delete the audio
	if err := h.audioService.DeleteAudio(audioID); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to delete audio"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, response.MessageResponse{Message: "Audio deleted successfully"})
}
