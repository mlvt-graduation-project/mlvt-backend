package mlvt_handler

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
	"mlvt/internal/service/video_service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MlvtController struct {
	audioService         audio_service.AudioService
	transcriptionService transcription_service.TranscriptionService
	videoService         video_service.VideoService
}

func NewMlvtController(audioService audio_service.AudioService, transcriptionService transcription_service.TranscriptionService, videoService video_service.VideoService) *MlvtController {
	return &MlvtController{
		audioService:         audioService,
		transcriptionService: transcriptionService,
		videoService:         videoService,
	}
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
func (h *MlvtController) ProcessTextToSpeech(c *gin.Context) {
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

// ProcessSpeechToText godoc
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
func (h *MlvtController) ProcessSpeechToText(c *gin.Context) {
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
			BaseRequest: request.BaseRequest{
				InputFileName:  videoFileName,
				InputLink:      videoDownloadURL,
				OutputFileName: transcriptionFileName,
				OutputLink:     transcriptionUploadURL,
				Model:          "",
			},
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
		var ec2Response response.EC2Response
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

// ProcessTextToText handles the translation of a transcription asynchronously
// @Summary Process transcription to translation asynchronously
// @Description Translates a transcription by processing it through an external service asynchronously
// @Tags transcriptions
// @Accept json
// @Produce json
// @Param transcription_id path uint64 true "ID of the transcription to translate"
// @Param source_language query string true "Source language code"
// @Param target_language query string true "Target language code"
// @Success 202 {object} response.MessageCreateResponseWithID "Accepted for processing"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transcriptions/translate/{transcription_id} [post]
func (h *MlvtController) ProcessTextToText(c *gin.Context) {
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
	transcription, _, err := h.transcriptionService.GetTranscriptionByID(transcriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "internal server error"})
		return
	}
	if transcription == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "transcription not found"})
		return
	}

	// Validate source language
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

	// Get folder from env config or use a predefined folder
	folder := env.EnvConfig.TranscriptionsFolder
	if folder == "" {
		folder = "transcriptions"
	}

	// Generate unique file name for the translated transcription
	translatedFileName := fmt.Sprintf("transcription_%d_%s.txt", transcriptionID, targetLang)

	// Create a new Transcription entity with status 'processing' and store in DB
	newTranscription := &entity.Transcription{
		VideoID:   transcription.VideoID,
		UserID:    transcription.UserID,
		Lang:      targetLang,
		Folder:    folder,
		FileName:  translatedFileName,
		Status:    entity.StatusProcessing,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	translatedTranscriptionID, err := h.transcriptionService.CreateTranscription(newTranscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to store translated transcription data"})
		return
	}

	// Respond to the frontend immediately with the new transcription and status 'processing'
	c.JSON(http.StatusAccepted, response.MessageCreateResponseWithID{
		Message: "Accepted for processing",
		Id:      translatedTranscriptionID,
	})

	// Start the asynchronous processing in a separate goroutine
	go func(originalTranscription *entity.Transcription, translatedTranscription *entity.Transcription, translatedTranscriptionID uint64, folder string) {
		// Ensure any panic in the goroutine does not crash the application
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered in goroutine: %v\n", r)
				h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			}
		}()

		// Generate presigned download URL for original transcription
		originalDownloadURL, err := h.transcriptionService.GeneratePresignedDownloadURL(originalTranscription.ID)
		if err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to generate original transcription download URL for transcription ID %d: %v\n", translatedTranscriptionID, err)
			return
		}

		// Generate presigned upload URL for translated transcription
		fileType := "text/plain"
		translationUploadURL, err := h.transcriptionService.GeneratePresignedUploadURL(folder, translatedTranscription.FileName, fileType)
		if err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to generate translation upload URL for transcription ID %d: %v\n", translatedTranscriptionID, err)
			return
		}

		// Prepare request payload matching TTTRequest struct
		requestPayload := request.TTTRequest{
			BaseRequest: request.BaseRequest{
				InputFileName:  originalTranscription.FileName,
				InputLink:      originalDownloadURL,
				OutputFileName: translatedTranscription.FileName,
				OutputLink:     translationUploadURL,
				Model:          "",
			},
			BaseLang: request.BaseLang{
				SourceLang: sourceLang,
				TargetLang: targetLang,
			},
		}

		jsonData, err := json.Marshal(requestPayload)
		if err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to marshal request payload for transcription ID %d: %v\n", translatedTranscriptionID, err)
			return
		}

		// Create a custom HTTP client with a timeout
		client := &http.Client{
			Timeout: 5 * time.Minute, // Must match EC2's handler timeout
		}

		// Send request to EC2 server
		ec2ServerURL := fmt.Sprintf("http://%s:%s/ttt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		resp, err := client.Post(ec2ServerURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to send request to EC2 server for transcription ID %d: %v\n", translatedTranscriptionID, err)
			return
		}
		defer resp.Body.Close()

		// Read and parse the EC2 response
		var ec2Response response.EC2Response
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to read EC2 response for transcription ID %d: %v\n", translatedTranscriptionID, err)
			return
		}

		if err := json.Unmarshal(bodyBytes, &ec2Response); err != nil {
			// Update transcription status to 'failed'
			if updateErr := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed); updateErr != nil {
				fmt.Printf("Failed to update transcription status: %v\n", updateErr)
			}
			fmt.Printf("Failed to parse EC2 response for transcription ID %d: %v\n", translatedTranscriptionID, err)
			return
		}

		// Handle EC2 response based on status
		switch ec2Response.Status {
		case "succeeded":
			// Update the transcription with the received text and set the status to 'succeeded'
			updateTranscription := &entity.Transcription{
				ID:        translatedTranscriptionID,
				Text:      ec2Response.Result,
				Status:    entity.StatusSucceeded,
				UpdatedAt: time.Now(),
			}

			if err := h.transcriptionService.UpdateTranscription(updateTranscription); err != nil {
				fmt.Printf("Failed to update transcription data for transcription ID %d: %v\n", translatedTranscriptionID, err)
				return
			}

			fmt.Printf("Successfully processed translation transcription ID %d\n", translatedTranscriptionID)

		case "failed":
			// Update transcription status to 'failed' with error message
			if err := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed); err != nil {
				fmt.Printf("Failed to update transcription status for transcription ID %d: %v\n", translatedTranscriptionID, err)
			}
			fmt.Printf("EC2 TTT processing failed for transcription ID %d: %s\n", translatedTranscriptionID, ec2Response.Error)

		case "timeout":
			// Update transcription status to 'failed' (or a specific 'timeout' status if defined)
			if err := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed); err != nil {
				fmt.Printf("Failed to update transcription status for transcription ID %d: %v\n", translatedTranscriptionID, err)
			}
			fmt.Printf("EC2 TTT processing timed out for transcription ID %d\n", translatedTranscriptionID)

		default:
			// Handle unexpected status
			if err := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed); err != nil {
				fmt.Printf("Failed to update transcription status for transcription ID %d: %v\n", translatedTranscriptionID, err)
			}
			fmt.Printf("EC2 TTT processing returned unknown status '%s' for transcription ID %d\n", ec2Response.Status, translatedTranscriptionID)
		}
	}(transcription, newTranscription, translatedTranscriptionID, folder)
}
