package transcription_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/request"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/transcription_service"
	"mlvt/internal/service/video_service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TranscriptionController struct {
	transcriptionService transcription_service.TranscriptionService
	videoService         video_service.VideoService
}

func NewTranscriptionController(
	transcriptionService transcription_service.TranscriptionService,
	videoService video_service.VideoService,
) *TranscriptionController {
	return &TranscriptionController{
		transcriptionService: transcriptionService,
		videoService:         videoService,
	}
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
	log.Warnf("extract from param: %s", transcriptionIDStr)
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

	id, err := h.transcriptionService.CreateTranscription(&transcription)
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

// ProcessVideoToTranscription godoc
// @Summary Process video to transcription asynchronously
// @Description Converts a video to transcription by processing it through an external service asynchronously
// @Tags transcriptions
// @Accept json
// @Produce json
// @Param video_id path uint64 true "ID of the video to process"
// @Success 202 {object} response.TranscriptionResponse "Accepted for processing"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transcriptions/process/{video_id} [post]
func (h *TranscriptionController) ProcessVideoToTranscription(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid video ID"})
		return
	}

	// Check if video exists
	video, _, _, err := h.videoService.GetVideoByID(videoID)
	if err != nil {
		if err.Error() == "video not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "video not found"})
		} else {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	// Get folder from environment config or use a predefined folder
	folder := env.EnvConfig.TranscriptionsFolder
	if folder == "" {
		folder = "transcriptions"
	}

	// Generate unique file name for transcription
	transcriptionFileName := fmt.Sprintf("transcription_%d.txt", videoID)

	// Create Transcription entity with status 'processing' and store in DB
	transcription := &entity.Transcription{
		VideoID:   videoID,
		UserID:    video.UserID,
		Folder:    folder,
		FileName:  transcriptionFileName,
		Status:    entity.StatusProcessing,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	transcriptionID, err := h.transcriptionService.CreateTranscription(transcription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to store transcription data"})
		return
	}

	// Respond to the frontend immediately with the new transcription and status 'processing'
	c.JSON(http.StatusAccepted, response.MessageCreateResponseWithID{
		Message: "Accepted for processing",
		Id:      transcriptionID,
	})

	// Start the asynchronous processing in a separate goroutine
	go func(transcriptionID uint64, videoID uint64, videoFileName string, folder string) {
		// Ensure any panic in the goroutine does not crash the application
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered in goroutine: %v\n", r)
				h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed)
			}
		}()

		// Generate presigned download URL for video
		videoDownloadURL, err := h.videoService.GeneratePresignedDownloadURLForVideo(videoID)
		if err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to generate video download URL for transcription ID %d: %v\n", transcriptionID, err)
			return
		}

		// Generate presigned upload URL for transcription
		fileType := "text/plain"
		transcriptionUploadURL, err := h.transcriptionService.GeneratePresignedUploadURL(folder, transcriptionFileName, fileType)
		if err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to generate transcription upload URL for transcription ID %d: %v\n", transcriptionID, err)
			return
		}

		// Prepare request payload matching STTRequest struct
		requestPayload := request.STTRequest{
			InputFileName:  videoFileName,
			InputLink:      videoDownloadURL,
			OutputFileName: transcriptionFileName,
			OutputLink:     transcriptionUploadURL,
		}

		jsonData, err := json.Marshal(requestPayload)
		if err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to marshal request payload for transcription ID %d: %v\n", transcriptionID, err)
			return
		}

		// Create a custom HTTP client with a timeout
		client := &http.Client{
			Timeout: 5 * time.Minute, // Must match EC2's handler timeout
		}

		// Send request to EC2 server
		ec2ServerURL := fmt.Sprintf("http://%s:%s/stt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		resp, err := client.Post(ec2ServerURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to send request to EC2 server for transcription ID %d: %v\n", transcriptionID, err)
			return
		}
		defer resp.Body.Close()

		// Read and parse the EC2 response
		var ec2Response response.EC2STTResponse
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to read EC2 response for transcription ID %d: %v\n", transcriptionID, err)
			return
		}

		if err := json.Unmarshal(bodyBytes, &ec2Response); err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to parse EC2 response for transcription ID %d: %v\n", transcriptionID, err)
			return
		}

		// Handle EC2 response based on status
		switch ec2Response.Status {
		case "succeeded":
			// Update the transcription with the received text and set the status to 'succeeded'
			updateTranscription := &entity.Transcription{
				ID:        transcriptionID,
				Text:      ec2Response.Result,
				Status:    entity.StatusSucceeded,
				UpdatedAt: time.Now(),
			}

			if err := h.transcriptionService.UpdateTranscription(updateTranscription); err != nil {
				fmt.Printf("Failed to update transcription data for transcription ID %d: %v\n", transcriptionID, err)
				return
			}

			fmt.Printf("Successfully processed transcription ID %d\n", transcriptionID)

		case "failed":
			// Update transcription status to 'failed' with error message
			if err := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed); err != nil {
				fmt.Printf("Failed to update transcription status for transcription ID %d: %v\n", transcriptionID, err)
			}
			fmt.Printf("EC2 STT processing failed for transcription ID %d: %s\n", transcriptionID, ec2Response.Error)

		case "timeout":
			// Update transcription status to 'timeout' or handle accordingly
			if err := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed); err != nil {
				fmt.Printf("Failed to update transcription status for transcription ID %d: %v\n", transcriptionID, err)
			}
			fmt.Printf("EC2 STT processing timed out for transcription ID %d\n", transcriptionID)

		default:
			// Handle unexpected status
			if err := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed); err != nil {
				fmt.Printf("Failed to update transcription status for transcription ID %d: %v\n", transcriptionID, err)
			}
			fmt.Printf("EC2 STT processing returned unknown status '%s' for transcription ID %d\n", ec2Response.Status, transcriptionID)
		}
	}(transcriptionID, videoID, video.FileName, folder)
}

// ProcessTranscriptionToTranslation handles the translation of a transcription
// @Summary Process transcription to translation
// @Description Translates a transcription by processing it through an external service
// @Tags transcriptions
// @Accept json
// @Produce json
// @Param transcription_id path uint64 true "ID of the transcription to translate"
// @Param source_language query string true "Source language code"
// @Param target_language query string true "Target language code"
// @Success 200 {object} response.TranscriptionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transcriptions/translate/{transcription_id} [post]
func (h *TranscriptionController) ProcessTranscriptionToTranslation(c *gin.Context) {
	transcriptionIDStr := c.Param("transcription_id")
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transcription ID"})
		return
	}

	sourceLang := c.Query("source_language")
	targetLang := c.Query("target_language")

	if sourceLang == "" || targetLang == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "source_language and target_language are required"})
		return
	}

	// Retrieve the existing transcription
	transcription, transcriptionDownloadURL, err := h.transcriptionService.GetTranscriptionByID(transcriptionID)
	if err != nil {
		if err.Error() == "transcription not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "transcription not found"})
		} else {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	// Check if transcription.Lang is empty or equals sourceLang
	if transcription.Lang == "" {
		// Update transcription.Lang to sourceLang
		transcription.Lang = sourceLang
		err = h.transcriptionService.UpdateTranscription(transcription)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to update transcription language"})
			return
		}
	} else if transcription.Lang != sourceLang {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "source_language does not match transcription language"})
		return
	}

	// Generate unique file name for the translated transcription
	translatedFileName := fmt.Sprintf("transcription_%d_%s.txt", transcriptionID, targetLang)

	// Get folder from env config or use a predefined folder
	folder := env.EnvConfig.TranscriptionsFolder
	if folder == "" {
		folder = "transcriptions"
	}

	// Generate presigned upload URL for the translated transcription
	fileType := "text/plain"
	translationUploadURL, err := h.transcriptionService.GeneratePresignedUploadURL(folder, translatedFileName, fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to generate translation upload URL"})
		return
	}

	// Prepare request payload matching TTTRequest struct
	requestPayload := request.TTTRequest{
		InputFileName:  transcription.FileName,
		InputLink:      transcriptionDownloadURL,
		OutputFileName: translatedFileName,
		OutputLink:     translationUploadURL,
		SourceLang:     sourceLang,
		TargetLang:     targetLang,
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to marshal request body"})
		return
	}

	// Send request to EC2 server
	ec2ServerURL := fmt.Sprintf("http://%s:%s/ttt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
	req, err := http.NewRequest("POST", ec2ServerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to create request to processing server"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to send request to processing server"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: fmt.Sprintf("processing server returned error: %s", string(respBody))})
		return
	}

	// Create new Transcription entity with targetLang
	newTranscription := &entity.Transcription{
		VideoID:  transcription.VideoID,
		UserID:   transcription.UserID,
		Lang:     targetLang,
		Folder:   folder,
		FileName: translatedFileName,
	}

	_, err = h.transcriptionService.CreateTranscription(newTranscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to store translated transcription data"})
		return
	}

	// Return the new transcription info to frontend
	c.JSON(http.StatusOK, response.TranscriptionResponse{
		Transcription: *newTranscription,
	})
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
func (h *TranscriptionController) UpdateTranscriptionStatus(c *gin.Context) {
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

	err = h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, req.Status)
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
