package audio_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/pkg/request"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/audio_service"
	"mlvt/internal/service/transcription_service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AudioController struct {
	audioService         audio_service.AudioService
	transcriptionService transcription_service.TranscriptionService
}

func NewAudioController(audioService audio_service.AudioService, transcriptionService transcription_service.TranscriptionService) *AudioController {
	return &AudioController{
		audioService:         audioService,
		transcriptionService: transcriptionService,
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

	// Create the audio and retrieve its ID
	audioID, err := h.audioService.CreateAudio(&audio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.MessageCreateResponseWithID{
		Message: "Audio added successfully",
		Id:      audioID,
	})
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
// @Param audio_id path uint64 true "ID of the audio file"
// @Param user_id path uint64 true "ID of the user"
// @Success 200 {object} response.AudioResponse "audio, download_url"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /audios/{audio_id}/user/{user_id} [get]
func (h *AudioController) GetAudioByUser(c *gin.Context) {
	// Parse audio ID from the URL path
	audioIDStr := c.Param("audio_id")
	audioID, err := strconv.ParseUint(audioIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid audio ID"})
		return
	}

	// Parse user ID from the URL path
	userIDStr := c.Param("user_id")
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
// @Param audio_id path uint64 true "ID of the audio file"
// @Param video_id path uint64 true "ID of the video"
// @Success 200 {object} response.AudioResponse "audio, download_url"
// @Failure 400 {object} response.ErrorResponse "error"
// @Failure 500 {object} response.ErrorResponse "error"
// @Router /audios/{audio_id}/video/{video_id} [get]
func (h *AudioController) GetAudioByVideoID(c *gin.Context) {
	// Parse audio ID from the URL path
	audioIDStr := c.Param("audio_id")
	audioID, err := strconv.ParseUint(audioIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid audio ID"})
		return
	}

	// Parse video ID from the URL path
	videoIDStr := c.Param("video_id")
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

// ProcessTextToSpeech godoc
// @Summary Convert transcription to speech asynchronously
// @Description Converts a transcription to speech by processing it through an external service asynchronously
// @Tags audios
// @Accept json
// @Produce json
// @Param transcription_id path uint64 true "ID of the transcription to convert"
// @Success 202 {object} response.MessageCreateResponseWithID "Accepted for processing"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /audios/process/{transcription_id} [post]
func (h *AudioController) ProcessTextToSpeech(c *gin.Context) {
	transcriptionIDStr := c.Param("transcription_id")
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transcription ID"})
		return
	}

	// // Optional: Authentication and Authorization
	// userID, exists := c.Get("user_id")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "unauthorized"})
	// 	return
	// }

	// Retrieve the existing transcription
	transcription, _, err := h.transcriptionService.GetTranscriptionByID(transcriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		return
	}
	if transcription == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "transcription not found"})
		return
	}

	// // Verify ownership
	// if transcription.UserID != userID.(uint64) {
	// 	c.JSON(http.StatusForbidden, response.ErrorResponse{Error: "forbidden"})
	// 	return
	// }

	// Get folder from environment config or use a predefined folder
	folder := env.EnvConfig.AudioFolder
	if folder == "" {
		folder = "audios"
	}

	// Generate unique file name for audio
	audioFileName := fmt.Sprintf("audio_%d.mp3", transcriptionID)

	// Create Audio entity with status 'processing' and store in DB
	audio := &entity.Audio{
		VideoID:   transcription.VideoID,
		UserID:    transcription.UserID,
		Duration:  0, // Initialize with 0; update later if needed
		Lang:      transcription.Lang,
		Folder:    folder,
		FileName:  audioFileName,
		Status:    entity.StatusProcessing,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	audioID, err := h.audioService.CreateAudio(audio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to store audio data"})
		return
	}

	// Respond to the frontend immediately with the new audio ID and status 'processing'
	c.JSON(http.StatusAccepted, response.MessageCreateResponseWithID{
		Message: "Accepted for processing",
		Id:      audioID,
	})

	// Start the asynchronous processing in a separate goroutine
	go func(audioID uint64, transcriptionID uint64, transcriptionLang string, folder string, transcriptionFileName string) {
		// Ensure any panic in the goroutine does not crash the application
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in goroutine: %v", r)
				if err := h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed); err != nil {
					log.Printf("Failed to update audio status after panic: %v", err)
				}
			}
		}()

		// Generate presigned download URL for transcription text
		transcriptionDownloadURL, err := h.transcriptionService.GeneratePresignedDownloadURL(transcriptionID)
		if err != nil {
			// Update audio status to 'failed'
			if updateErr := h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed); updateErr != nil {
				log.Printf("Failed to update audio status: %v", updateErr)
			}
			log.Printf("Failed to generate transcription download URL for audio ID %d: %v", audioID, err)
			return
		}

		// Generate presigned upload URL for audio
		fileType := "audio/mpeg"
		audioUploadURL, err := h.audioService.GeneratePresignedUploadURL(folder, audioFileName, fileType)
		if err != nil {
			// Update audio status to 'failed'
			if updateErr := h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed); updateErr != nil {
				log.Printf("Failed to update audio status: %v", updateErr)
			}
			log.Printf("Failed to generate audio upload URL for audio ID %d: %v", audioID, err)
			return
		}

		// Prepare request payload matching TTSRequest struct
		requestPayload := request.TTSRequest{
			BaseRequest: request.BaseRequest{
				InputFileName:  transcriptionFileName, // Use the transcription file name
				InputLink:      transcriptionDownloadURL,
				OutputFileName: audioFileName,
				OutputLink:     audioUploadURL,
				Model:          "",
			},
			Lang: transcriptionLang,
		}

		jsonData, err := json.Marshal(requestPayload)
		if err != nil {
			// Update audio status to 'failed'
			if updateErr := h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed); updateErr != nil {
				log.Printf("Failed to update audio status: %v", updateErr)
			}
			log.Printf("Failed to marshal request payload for audio ID %d: %v", audioID, err)
			return
		}

		// Create a custom HTTP client with a timeout
		client := &http.Client{
			Timeout: 5 * time.Minute, // Must match EC2's handler timeout
		}

		// Send request to EC2 server
		ec2ServerURL := fmt.Sprintf("http://%s:%s/tts", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		resp, err := client.Post(ec2ServerURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			// Update audio status to 'failed'
			if updateErr := h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed); updateErr != nil {
				log.Printf("Failed to update audio status: %v", updateErr)
			}
			log.Printf("Failed to send request to EC2 server for audio ID %d: %v", audioID, err)
			return
		}
		defer resp.Body.Close()

		// Read and parse the EC2 response
		var ec2Response response.EC2Response
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// Update audio status to 'failed'
			if updateErr := h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed); updateErr != nil {
				log.Printf("Failed to update audio status: %v", updateErr)
			}
			log.Printf("Failed to read EC2 response for audio ID %d: %v", audioID, err)
			return
		}

		if err := json.Unmarshal(bodyBytes, &ec2Response); err != nil {
			// Update audio status to 'failed'
			if updateErr := h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed); updateErr != nil {
				log.Printf("Failed to update audio status: %v", updateErr)
			}
			log.Printf("Failed to parse EC2 response for audio ID %d: %v", audioID, err)
			return
		}

		// Handle EC2 response based on status
		switch ec2Response.Status {
		case "succeeded":
			updateAudio := &entity.Audio{
				ID:        audioID,
				Status:    entity.StatusSucceeded,
				UpdatedAt: time.Now(),
			}

			// if ec2Response.Duration > 0 {
			// 	updateAudio.Duration = ec2Response.Duration
			// }

			if err := h.audioService.UpdateAudio(updateAudio); err != nil {
				log.Printf("Failed to update audio data for audio ID %d: %v", audioID, err)
				return
			}

			log.Printf("Successfully processed audio ID %d", audioID)

		case "failed":
			// Update audio status to 'failed' with error message
			if err := h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed); err != nil {
				log.Printf("Failed to update audio status for audio ID %d: %v", audioID, err)
			}
			log.Printf("EC2 TTS processing failed for audio ID %d: %s", audioID, ec2Response.Error)

		case "timeout":
			// Update audio status to 'failed' or a specific 'timeout' status if defined
			if err := h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed); err != nil {
				log.Printf("Failed to update audio status for audio ID %d: %v", audioID, err)
			}
			log.Printf("EC2 TTS processing timed out for audio ID %d", audioID)

		default:
			// Handle unexpected status
			if err := h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed); err != nil {
				log.Printf("Failed to update audio status for audio ID %d: %v", audioID, err)
			}
			log.Printf("EC2 TTS processing returned unknown status '%s' for audio ID %d", ec2Response.Status, audioID)
		}
	}(audioID, transcriptionID, transcription.Lang, folder, transcription.FileName)
}
