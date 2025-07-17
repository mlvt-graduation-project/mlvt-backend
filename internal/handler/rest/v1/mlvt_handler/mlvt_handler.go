package mlvt_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/request"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/media_service"
	"mlvt/internal/service/notify_service"
	"mlvt/internal/service/progress_service"
	"mlvt/internal/service/traffic_service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MlvtController struct {
	mediaService    media_service.MediaService
	progressService progress_service.ProgressService
	trafficService  traffic_service.TrafficService
	notifyService   notify_service.NotifyService
}

func NewMlvtController(
	mediaService media_service.MediaService,
	progressService progress_service.ProgressService,
	trafficService traffic_service.TrafficService,
	notifyService notify_service.NotifyService,
) *MlvtController {
	return &MlvtController{
		mediaService:    mediaService,
		progressService: progressService,
		trafficService:  trafficService,
		notifyService:   notifyService,
	}
}

// Helper function to send requests to EC2 and handle the response
func sendRequestToEC2(requestPayload interface{}, ec2Endpoint string, timeout time.Duration) (*response.EC2Response, error) {
	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Post(ec2Endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to EC2: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read EC2 response: %v", err)
	}

	var ec2Response response.EC2Response
	if err := json.Unmarshal(bodyBytes, &ec2Response); err != nil {
		return nil, fmt.Errorf("failed to parse EC2 response: %v", err)
	}

	return &ec2Response, nil
}

func (h *MlvtController) quickLogTraffic(trafficType entity.TrafficActionType, userID uint64, entityID uint64) {
	ctx := context.Background()
	switch trafficType {
	case entity.ProcessSTTModelAction:
		if _, err := h.trafficService.CreateTraffic(ctx, entity.Traffic{
			ActionType:     trafficType,
			Description:    fmt.Sprintf("process new STT pipeline, text id: %d", entityID),
			UserID:         userID,
			UserPermission: entity.UserRole,
			Timestamp:      time.Now().Unix(),
		}); err != nil {
			log.Errorf("failed to log traffic of %s", trafficType)
		}
		return
	case entity.ProcessTTTModelAction:
		if _, err := h.trafficService.CreateTraffic(ctx, entity.Traffic{
			ActionType:     trafficType,
			Description:    fmt.Sprintf("process new TTT pipeline, text id: %d", entityID),
			UserID:         userID,
			UserPermission: entity.UserRole,
			Timestamp:      time.Now().Unix(),
		}); err != nil {
			log.Errorf("failed to log traffic of %s", trafficType)
		}
		return
	case entity.ProcessTTSModelAction:
		if _, err := h.trafficService.CreateTraffic(ctx, entity.Traffic{
			ActionType:     trafficType,
			Description:    fmt.Sprintf("process new TTS pipeline, audio id: %d", entityID),
			UserID:         userID,
			UserPermission: entity.UserRole,
			Timestamp:      time.Now().Unix(),
		}); err != nil {
			log.Errorf("failed to log traffic of %s", trafficType)
		}
		return
	case entity.ProcessLSModelAction:
		if _, err := h.trafficService.CreateTraffic(ctx, entity.Traffic{
			ActionType:     trafficType,
			Description:    fmt.Sprintf("process new LS pipeline, video id: %d", entityID),
			UserID:         userID,
			UserPermission: entity.UserRole,
			Timestamp:      time.Now().Unix(),
		}); err != nil {
			log.Errorf("failed to log traffic of %s", trafficType)
		}
		return
	case entity.ProcessFullPipelineModelAction:
		if _, err := h.trafficService.CreateTraffic(ctx, entity.Traffic{
			ActionType:     trafficType,
			Description:    fmt.Sprintf("process new Full pipeline, video id: %d", entityID),
			UserID:         userID,
			UserPermission: entity.UserRole,
			Timestamp:      time.Now().Unix(),
		}); err != nil {
			log.Errorf("failed to log traffic of %s", trafficType)
		}
		return
	default:
		return
	}
}

// Helper function to send notifications with MLVT header
func (h *MlvtController) sendNotification(message string) {
	go func() {
		// Add MLVT Backend header to all messages
		formattedMessage := h.formatNotificationMessage(message)
		if err := h.notifyService.SendNotification(formattedMessage); err != nil {
			log.Errorf("Failed to send notification: %v", err)
		}
	}()
}

// Helper function to format notification message with header
func (h *MlvtController) formatNotificationMessage(message string) string {
	appEnv := env.EnvConfig.AppEnv
	if appEnv == "" {
		appEnv = "unknown"
	}

	header := fmt.Sprintf("🤖 <b>MLVT Backend [%s]</b>\n"+
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n",
		strings.ToUpper(appEnv))

	footer := "\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		fmt.Sprintf("⏰ <i>%s</i>", time.Now().Format("2006-01-02 15:04:05"))

	return header + message + footer
}

// ProcessSpeechToText godoc
// @Summary Convert video to transcription asynchronously
// @Description Converts a video to text using speech-to-text processing asynchronously
// @Tags transcriptions
// @Accept  json
// @Produce  json
// @Param   video_id   path    uint64     true  "Video ID"
// @Success 202 {object} response.MessageCreateResponseWithID "Accepted for processing"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transcriptions/process/{video_id} [post]
func (h *MlvtController) ProcessSpeechToText(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid video ID"})
		return
	}
	sourceLang := c.Query("source_language")
	if sourceLang == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "source_language and target_language are required"})
		return
	}

	video, _, _, err := h.mediaService.GetVideoByID(videoID)
	if err != nil || video == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Video not found"})
		return
	}

	folder := env.EnvConfig.TranscriptionsFolder
	if folder == "" {
		folder = "transcriptions"
	}

	transcriptionFileName := fmt.Sprintf("transcription_%d.txt", videoID)
	transcription := &entity.Transcription{
		VideoID:   videoID,
		UserID:    video.UserID,
		Folder:    folder,
		Lang:      sourceLang,
		FileName:  transcriptionFileName,
		Status:    entity.StatusProcessing,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	transcriptionID, err := h.mediaService.CreateTranscription(transcription, false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store transcription data"})
		return
	}

	h.quickLogTraffic(entity.ProcessSTTModelAction, video.UserID, transcriptionID)

	// Send notification about STT processing start
	h.sendNotification(fmt.Sprintf("🎤 <b>Speech-to-Text Processing Started</b>\n\n"+
		"📹 Video ID: %d\n"+
		"👤 User ID: %d\n"+
		"📝 Transcription ID: %d\n"+
		"🗣️ Language: %s\n\n"+
		"⏳ Processing in progress...", videoID, video.UserID, transcriptionID, sourceLang))

	// Insert to mongodb
	sttDocument := &entity.Progress{
		UserID:                    video.UserID,
		ProgressType:              entity.ProgressTypeSTT,
		OriginalVideoID:           videoID,
		OriginalTranscriptionID:   transcriptionID,
		TranslatedTranscriptionID: 0,
		AudioID:                   0,
		ProgressedVideoID:         0,
		Status:                    entity.StatusProcessing,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}

	documentId, err := h.progressService.Create(context.Background(), *sttDocument)
	if err != nil {
		log.Errorf("Failed to insert document ", err)
	}
	log.Infof("Added document STT, Id: ", documentId)

	// Respond immediately to the client
	c.JSON(http.StatusAccepted, response.MessageCreateResponseWithID{
		Message: "Accepted for processing",
		Id:      transcriptionID,
	})

	// Start asynchronous processing
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Warnf("Recovered in goroutine: %v", r)
				h.mediaService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed)
				h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			}
		}()

		videoDownloadURL, err := h.mediaService.GeneratePresignedDownloadURLForVideo(videoID)
		if err != nil {
			h.mediaService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate video download URL: %v", err)
			return
		}

		fileType := "text/plain"
		transcriptionUploadURL, err := h.mediaService.GeneratePresignedUploadURLForText(folder, transcriptionFileName, fileType)
		if err != nil {
			h.mediaService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate transcription upload URL: %v", err)
			return
		}

		requestPayload := request.STTRequest{
			BaseRequest: request.BaseRequest{
				InputFileName:  video.FileName,
				InputLink:      videoDownloadURL,
				OutputFileName: transcriptionFileName,
				OutputLink:     transcriptionUploadURL,
				Model:          "",
			},
		}

		// Marshal the payload to JSON for curl logging
		payloadBytes, err := json.Marshal(requestPayload)
		if err != nil {
			log.Warnf("Error marshaling STT payload: %v", err)
		}

		ec2ServerURL := fmt.Sprintf("http://%s:%s/stt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)

		// Build and print the curl command for debugging
		curlCmd := fmt.Sprintf(`curl -X POST "%s" -H "Content-Type: application/json" -d '%s'`, ec2ServerURL, string(payloadBytes))
		fmt.Println("STT Curl Command:", curlCmd)

		ec2Response, err := sendRequestToEC2(requestPayload, ec2ServerURL, 5*time.Minute)
		log.Infof("ec2 response: %v\n\n", ec2Response)
		if err != nil || ec2Response.Status != "succeeded" {
			h.mediaService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("EC2 processing failed: %v", err)

			// Send failure notification
			h.sendNotification(fmt.Sprintf("❌ <b>Speech-to-Text Processing Failed</b>\n\n"+
				"📹 Video ID: %d\n"+
				"👤 User ID: %d\n"+
				"📝 Transcription ID: %d\n"+
				"🗣️ Language: %s\n\n"+
				"💥 Error: %v", videoID, video.UserID, transcriptionID, sourceLang, err))
			return
		}

		updateTranscription := &entity.Transcription{
			ID:        transcriptionID,
			Text:      ec2Response.Result,
			UpdatedAt: time.Now(),
		}

		if err := h.mediaService.UpdateTranscription(updateTranscription); err != nil {
			log.Errorf("Failed to update transcription data: %v", err)
		}

		// Update status to succeeded
		if err := h.mediaService.UpdateTranscriptionStatus(transcriptionID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update transcription status: %v", err)
		}

		h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusSucceeded)

		// Send success notification
		h.sendNotification(fmt.Sprintf("✅ <b>Speech-to-Text Processing Completed</b>\n\n"+
			"📹 Video ID: %d\n"+
			"👤 User ID: %d\n"+
			"📝 Transcription ID: %d\n"+
			"🗣️ Language: %s\n\n"+
			"📄 Result: %s", videoID, video.UserID, transcriptionID, sourceLang,
			func() string {
				if len(ec2Response.Result) > 100 {
					return ec2Response.Result[:100] + "..."
				}
				return ec2Response.Result
			}()))
	}()
}

// ProcessTextToText godoc
// @Summary Translate transcription asynchronously
// @Description Translates a transcription from source language to target language asynchronously
// @Tags transcriptions
// @Accept  json
// @Produce  json
// @Param   transcription_id  path    uint64     true  "Transcription ID"
// @Param   source_language   query   string     true  "Source language code"
// @Param   target_language   query   string     true  "Target language code"
// @Success 202 {object} response.MessageCreateResponseWithID "Accepted for processing"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transcriptions/translate/{transcription_id} [post]
func (h *MlvtController) ProcessTextToText(c *gin.Context) {
	transcriptionIDStr := c.Param("transcription_id")
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid transcription ID"})
		return
	}

	sourceLang := c.Query("source_language")
	targetLang := c.Query("target_language")

	if sourceLang == "" || targetLang == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "source_language and target_language are required"})
		return
	}

	originalTranscription, _, err := h.mediaService.GetTranscriptionByID(transcriptionID)
	if err != nil || originalTranscription == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Transcription not found"})
		return
	}

	folder := env.EnvConfig.TranscriptionsFolder
	if folder == "" {
		folder = "transcriptions"
	}

	translatedFileName := fmt.Sprintf("transcription_%d_%s.txt", transcriptionID, targetLang)
	newTranscription := &entity.Transcription{
		VideoID:                 originalTranscription.VideoID,
		UserID:                  originalTranscription.UserID,
		OriginalTranscriptionID: transcriptionID,
		Lang:                    targetLang,
		Folder:                  folder,
		FileName:                translatedFileName,
		Status:                  entity.StatusProcessing,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	translatedTranscriptionID, err := h.mediaService.CreateTranscription(newTranscription, false, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store translated transcription data"})
		return
	}

	h.quickLogTraffic(entity.ProcessTTTModelAction, originalTranscription.UserID, translatedTranscriptionID)

	// Send notification about TTT processing start
	h.sendNotification(fmt.Sprintf("🔄 <b>Text-to-Text Translation Started</b>\n\n"+
		"📝 Original Transcription ID: %d\n"+
		"📝 New Transcription ID: %d\n"+
		"👤 User ID: %d\n"+
		"🗣️ From: %s → %s\n\n"+
		"⏳ Translating...", transcriptionID, translatedTranscriptionID, originalTranscription.UserID, sourceLang, targetLang))

	// Insert to mongodb
	sttDocument := &entity.Progress{
		UserID:                    originalTranscription.UserID,
		ProgressType:              entity.ProgressTypeTTT,
		OriginalVideoID:           originalTranscription.VideoID,
		OriginalTranscriptionID:   originalTranscription.ID,
		TranslatedTranscriptionID: translatedTranscriptionID,
		AudioID:                   0,
		ProgressedVideoID:         0,
		Status:                    entity.StatusProcessing,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}

	documentId, err := h.progressService.Create(context.Background(), *sttDocument)
	if err != nil {
		log.Errorf("Failed to insert document ", err)
	}
	log.Infof("Added document STT, Id: ", documentId)

	// Respond immediately to the client
	c.JSON(http.StatusAccepted, response.MessageCreateResponseWithID{
		Message: "Accepted for processing",
		Id:      translatedTranscriptionID,
	})

	// Start asynchronous processing
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Warnf("Recovered in goroutine: %v", r)
				h.mediaService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
				h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			}
		}()

		originalDownloadURL, err := h.mediaService.GeneratePresignedDownloadURLForText(transcriptionID)
		if err != nil {
			h.mediaService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate original transcription download URL: %v", err)
			return
		}

		fileType := "text/plain"
		translationUploadURL, err := h.mediaService.GeneratePresignedUploadURLForText(folder, translatedFileName, fileType)
		if err != nil {
			h.mediaService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate translation upload URL: %v", err)
			return
		}

		requestPayload := request.TTTRequest{
			BaseRequest: request.BaseRequest{
				InputFileName:  originalTranscription.FileName,
				InputLink:      originalDownloadURL,
				OutputFileName: translatedFileName,
				OutputLink:     translationUploadURL,
				Model:          "",
			},
			BaseLang: request.BaseLang{
				SourceLang: sourceLang,
				TargetLang: targetLang,
			},
		}

		// Marshal the payload to JSON for curl logging
		payloadBytes, err := json.Marshal(requestPayload)
		if err != nil {
			log.Warnf("Error marshaling TTT payload: %v", err)
		}

		ec2ServerURL := fmt.Sprintf("http://%s:%s/ttt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)

		// Build and print the curl command for debugging
		curlCmd := fmt.Sprintf(`curl -X POST "%s" -H "Content-Type: application/json" -d '%s'`, ec2ServerURL, string(payloadBytes))
		fmt.Println("TTT Curl Command:", curlCmd)

		ec2Response, err := sendRequestToEC2(requestPayload, ec2ServerURL, 5*time.Minute)
		if err != nil || ec2Response.Status != "succeeded" {
			h.mediaService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("EC2 processing failed: %v", err)

			// Send failure notification
			h.sendNotification(fmt.Sprintf("❌ <b>Text-to-Text Translation Failed</b>\n\n"+
				"📝 Original Transcription ID: %d\n"+
				"📝 New Transcription ID: %d\n"+
				"👤 User ID: %d\n"+
				"🗣️ From: %s → %s\n\n"+
				"💥 Error: %v", transcriptionID, translatedTranscriptionID, originalTranscription.UserID, sourceLang, targetLang, err))
			return
		}

		updateTranscription := &entity.Transcription{
			ID:        translatedTranscriptionID,
			Text:      ec2Response.Result,
			UpdatedAt: time.Now(),
		}

		if err := h.mediaService.UpdateTranscription(updateTranscription); err != nil {
			log.Errorf("Failed to update transcription data: %v", err)
		}

		// Update status to succeeded
		if err := h.mediaService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update transcription status: %v", err)
		}

		h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusSucceeded)

		// Send success notification
		h.sendNotification(fmt.Sprintf("✅ <b>Text-to-Text Translation Completed</b>\n\n"+
			"📝 Original Transcription ID: %d\n"+
			"📝 New Transcription ID: %d\n"+
			"👤 User ID: %d\n"+
			"🗣️ From: %s → %s\n\n"+
			"📄 Translated Text: %s", transcriptionID, translatedTranscriptionID, originalTranscription.UserID, sourceLang, targetLang,
			func() string {
				if len(ec2Response.Result) > 100 {
					return ec2Response.Result[:100] + "..."
				}
				return ec2Response.Result
			}()))
	}()
}

// ProcessTextToSpeech godoc
// @Summary Convert transcription to speech asynchronously
// @Description Converts a transcription to audio using text-to-speech processing asynchronously
// @Tags audios
// @Accept  json
// @Produce  json
// @Param   transcription_id  path    uint64     true  "Transcription ID"
// @Success 202 {object} response.MessageCreateResponseWithID "Accepted for processing"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /audios/process/{transcription_id} [post]
func (h *MlvtController) ProcessTextToSpeech(c *gin.Context) {
	transcriptionIDStr := c.Param("transcription_id")
	transcriptionID, err := strconv.ParseUint(transcriptionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid transcription ID"})
		return
	}

	videoIDStr := c.Query("video_id")
	audioIDStr := c.Query("audio_id")

	var inputAudioLink string
	var inputAudioFileName string

	if videoIDStr != "" {
		videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid video ID"})
			return
		}
		video, _, _, err := h.mediaService.GetVideoByID(videoID)
		if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Video not found"})
			return
		}
		videoDownloadURL, err := h.mediaService.GeneratePresignedDownloadURLForVideo(videoID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to generate video download URL"})
			return
		}
		inputAudioLink = videoDownloadURL
		inputAudioFileName = video.FileName
	} else if audioIDStr != "" {
		audioID, err := strconv.ParseUint(audioIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid audio ID"})
			return
		}
		audio, _, err := h.mediaService.GetAudioByID(audioID)
		if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Audio not found"})
			return
		}
		audioDownloadURL, err := h.mediaService.GeneratePresignedDownloadURL(audioID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to generate audio download URL"})
			return
		}
		inputAudioLink = audioDownloadURL
		inputAudioFileName = audio.FileName
	} else {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid video or audio ID"})
		return
	}

	transcription, _, err := h.mediaService.GetTranscriptionByID(transcriptionID)
	if err != nil || transcription == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Transcription not found"})
		return
	}

	folder := env.EnvConfig.AudioFolder
	if folder == "" {
		folder = "audios"
	}

	audioFileName := fmt.Sprintf("audio_%d.mp3", transcriptionID)
	audio := &entity.Audio{
		TranscriptionID: transcriptionID,
		VideoID:         transcription.VideoID,
		UserID:          transcription.UserID,
		Lang:            transcription.Lang,
		Folder:          folder,
		FileName:        audioFileName,
		Status:          entity.StatusProcessing,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	audioID, err := h.mediaService.CreateAudio(audio, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store audio data"})
		return
	}

	h.quickLogTraffic(entity.ProcessTTSModelAction, transcription.UserID, audioID)

	// Send notification about TTS processing start
	h.sendNotification(fmt.Sprintf("🎵 <b>Text-to-Speech Processing Started</b>\n\n"+
		"📝 Transcription ID: %d\n"+
		"🎧 Audio ID: %d\n"+
		"👤 User ID: %d\n"+
		"🗣️ Language: %s\n\n"+
		"⏳ Converting text to audio...", transcriptionID, audioID, transcription.UserID, transcription.Lang))

	// Insert to mongodb
	sttDocument := &entity.Progress{
		UserID:                    transcription.UserID,
		ProgressType:              entity.ProgressTypeTTS,
		OriginalVideoID:           transcription.VideoID,
		OriginalTranscriptionID:   transcription.OriginalTranscriptionID,
		TranslatedTranscriptionID: transcription.ID,
		AudioID:                   audioID,
		ProgressedVideoID:         0,
		Status:                    entity.StatusProcessing,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}

	documentId, err := h.progressService.Create(context.Background(), *sttDocument)
	if err != nil {
		log.Errorf("Failed to insert document ", err)
	}
	log.Infof("Added document STT, Id: ", documentId)

	// Respond immediately to the client
	c.JSON(http.StatusAccepted, response.MessageCreateResponseWithID{
		Message: "Accepted for processing",
		Id:      audioID,
	})

	// Start asynchronous processing
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Warnf("Recovered in goroutine: %v", r)
				h.mediaService.UpdateAudioStatus(audioID, entity.StatusFailed)
				h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			}
		}()

		transcriptionDownloadURL, err := h.mediaService.GeneratePresignedDownloadURLForText(transcriptionID)
		if err != nil {
			h.mediaService.UpdateAudioStatus(audioID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate transcription download URL: %v", err)
			return
		}

		fileType := "audio/mpeg"
		audioUploadURL, err := h.mediaService.GeneratePresignedUploadURL(folder, audioFileName, fileType)
		if err != nil {
			h.mediaService.UpdateAudioStatus(audioID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate audio upload URL: %v", err)
			return
		}

		requestPayload := request.TTSRequest{
			BaseRequest: request.BaseRequest{
				InputFileName:  transcription.FileName,
				InputLink:      transcriptionDownloadURL,
				OutputFileName: audioFileName,
				OutputLink:     audioUploadURL,
				Model:          "",
			},
			InputAudioFileName: inputAudioFileName,
			InputAudioLink:     inputAudioLink,
			Lang:               transcription.Lang,
		}

		// Marshal the payload to JSON for curl logging
		payloadBytes, err := json.Marshal(requestPayload)
		if err != nil {
			log.Warnf("Error marshaling TTS payload: %v", err)
		}

		ec2ServerURL := fmt.Sprintf("http://%s:%s/tts", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)

		// Build and print the curl command for debugging
		curlCmd := fmt.Sprintf(`curl -X POST "%s" -H "Content-Type: application/json" -d '%s'`, ec2ServerURL, string(payloadBytes))
		fmt.Println("TTS Curl Command:", curlCmd)

		ec2Response, err := sendRequestToEC2(requestPayload, ec2ServerURL, 5*time.Minute)
		if err != nil || ec2Response.Status != "succeeded" {
			h.mediaService.UpdateAudioStatus(audioID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("EC2 processing failed: %v", err)

			// Send failure notification
			h.sendNotification(fmt.Sprintf("❌ <b>Text-to-Speech Processing Failed</b>\n\n"+
				"📝 Transcription ID: %d\n"+
				"🎧 Audio ID: %d\n"+
				"👤 User ID: %d\n"+
				"🗣️ Language: %s\n\n"+
				"💥 Error: %v", transcriptionID, audioID, transcription.UserID, transcription.Lang, err))
			return
		}

		updateAudio := &entity.Audio{
			ID:        audioID,
			UpdatedAt: time.Now(),
		}

		if err := h.mediaService.UpdateAudio(updateAudio); err != nil {
			log.Errorf("Failed to update audio data: %v", err)
		}
		// Update status to succeeded
		if err := h.mediaService.UpdateAudioStatus(audioID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update audio status: %v", err)
		}

		h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusSucceeded)

		// Send success notification
		h.sendNotification(fmt.Sprintf("✅ <b>Text-to-Speech Processing Completed</b>\n\n"+
			"📝 Transcription ID: %d\n"+
			"🎧 Audio ID: %d\n"+
			"👤 User ID: %d\n"+
			"🗣️ Language: %s\n"+
			"📄 File: %s\n\n"+
			"🎉 Audio file generated successfully!", transcriptionID, audioID, transcription.UserID, transcription.Lang, audioFileName))
	}()
}

// ProcessLipSync godoc
// @Summary Perform lip synchronization asynchronously
// @Description Synchronizes lip movements in a video based on an audio track asynchronously
// @Tags lipsync
// @Accept  json
// @Produce  json
// @Param   video_id   path    uint64     true  "Video ID"
// @Param   audio_id   path    uint64     true  "Audio ID"
// @Success 202 {object} response.MessageCreateResponseWithID "Accepted for processing"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /lipsync/{video_id}/{audio_id} [post]
func (h *MlvtController) ProcessLipSync(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	audioIDStr := c.Param("audio_id")

	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid video ID"})
		return
	}

	audioID, err := strconv.ParseUint(audioIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid audio ID"})
		return
	}

	video, _, _, err := h.mediaService.GetVideoByID(videoID)
	if err != nil || video == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Video not found"})
		return
	}

	audio, _, err := h.mediaService.GetAudioByID(audioID)
	if err != nil || audio == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Audio not found"})
		return
	}

	folder := env.EnvConfig.VideosFolder
	if folder == "" {
		folder = "raw_videos"
	}

	outputVideoFileName := fmt.Sprintf("lipsync_%d_%d.mp4", videoID, audioID)
	outputVideo := &entity.Video{
		OriginalVideoID: videoID,
		AudioID:         audioID,
		UserID:          video.UserID,
		Folder:          folder,
		Image:           video.Image,
		FileName:        outputVideoFileName,
		Status:          entity.StatusProcessing,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	outputVideoID, err := h.mediaService.CreateVideo(outputVideo, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store output video data"})
		return
	}

	h.quickLogTraffic(entity.ProcessLSModelAction, video.UserID, outputVideoID)

	// Send notification about LipSync processing start
	h.sendNotification(fmt.Sprintf("👄 <b>Lip Sync Processing Started</b>\n\n"+
		"📹 Input Video ID: %d\n"+
		"🎧 Audio ID: %d\n"+
		"📹 Output Video ID: %d\n"+
		"👤 User ID: %d\n\n"+
		"⏳ Synchronizing lip movements...", videoID, audioID, outputVideoID, video.UserID))

	// Insert to mongodb
	sttDocument := &entity.Progress{
		UserID:                    video.UserID,
		ProgressType:              entity.ProgressTypeLS,
		OriginalVideoID:           videoID,
		TranslatedTranscriptionID: audio.TranscriptionID,
		AudioID:                   audioID,
		ProgressedVideoID:         outputVideoID,
		Status:                    entity.StatusProcessing,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}

	if audio.TranscriptionID != 0 {
		translatedTranscription, _, err := h.mediaService.GetTranscriptionByID(audio.TranscriptionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to get translated transcription by id"})
		}
		if translatedTranscription != nil {
			sttDocument.OriginalTranscriptionID = translatedTranscription.OriginalTranscriptionID
		}
	}

	documentId, err := h.progressService.Create(context.Background(), *sttDocument)
	if err != nil {
		log.Errorf("Failed to insert document ", err)
	}
	log.Infof("Added document STT, Id: ", documentId)

	// Respond immediately to the client
	c.JSON(http.StatusAccepted, response.MessageCreateResponseWithID{
		Message: "Accepted for processing",
		Id:      outputVideoID,
	})

	// Start asynchronous processing
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Warnf("Recovered in goroutine: %v", r)
				h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
				h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			}
		}()

		videoDownloadURL, err := h.mediaService.GeneratePresignedDownloadURLForVideo(videoID)
		if err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate video download URL: %v", err)
			return
		}

		audioDownloadURL, err := h.mediaService.GeneratePresignedDownloadURL(audioID)
		if err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate audio download URL: %v", err)
			return
		}

		fileType := "video/mp4"
		outputVideoUploadURL, err := h.mediaService.GeneratePresignedUploadURLForVideo(folder, outputVideoFileName, fileType)
		if err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate output video upload URL: %v", err)
			return
		}

		requestPayload := request.LSRequest{
			InputVideoFileName:  video.FileName,
			InputVideoLink:      videoDownloadURL,
			InputAudioFileName:  audio.FileName,
			InputAudioLink:      audioDownloadURL,
			OutputVideoFileName: outputVideoFileName,
			OutputVideoLink:     outputVideoUploadURL,
			Model:               "",
		}

		// Marshal the payload to JSON for curl logging
		payloadBytes, err := json.Marshal(requestPayload)
		if err != nil {
			log.Warnf("Error marshaling LS payload: %v", err)
		}

		ec2ServerURL := fmt.Sprintf("http://%s:%s/ls", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)

		// Build and print the curl command for debugging
		curlCmd := fmt.Sprintf(`curl -X POST "%s" -H "Content-Type: application/json" -d '%s'`, ec2ServerURL, string(payloadBytes))
		fmt.Println("LS Curl Command:", curlCmd)

		ec2Response, err := sendRequestToEC2(requestPayload, ec2ServerURL, 25*time.Minute)
		if err != nil || ec2Response.Status != "succeeded" {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("EC2 processing failed: %v", err)

			// Send failure notification
			h.sendNotification(fmt.Sprintf("❌ <b>Lip Sync Processing Failed</b>\n\n"+
				"📹 Input Video ID: %d\n"+
				"🎧 Audio ID: %d\n"+
				"📹 Output Video ID: %d\n"+
				"👤 User ID: %d\n\n"+
				"💥 Error: %v", videoID, audioID, outputVideoID, video.UserID, err))
			return
		}

		updateVideo := &entity.Video{
			ID:        outputVideoID,
			UpdatedAt: time.Now(),
		}

		if err := h.mediaService.UpdateVideo(updateVideo); err != nil {
			log.Errorf("Failed to update video data: %v", err)
		}

		// Update status to succeeded
		if err := h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update video status: %v", err)
		}

		h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusSucceeded)

		// Send success notification
		h.sendNotification(fmt.Sprintf("✅ <b>Lip Sync Processing Completed</b>\n\n"+
			"📹 Input Video ID: %d\n"+
			"🎧 Audio ID: %d\n"+
			"📹 Output Video ID: %d\n"+
			"👤 User ID: %d\n"+
			"📄 Output File: %s\n\n"+
			"🎉 Lip sync video generated successfully!", videoID, audioID, outputVideoID, video.UserID, outputVideoFileName))
	}()
}

// ProcessFullPipeline godoc
// @Summary Process full pipeline asynchronously
// @Description Processes a video through the full pipeline: Speech-to-Text, Text-to-Text, Text-to-Speech, and Lip Sync asynchronously
// @Tags pipeline
// @Accept  json
// @Produce  json
// @Param   video_id         path    uint64     true  "Video ID"
// @Param   source_language  query   string     true  "Source language code"
// @Param   target_language  query   string     true  "Target language code"
// @Success 202 {object} response.MessageCreateResponseWithID "Accepted for processing"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /pipeline/full/{video_id} [post]
func (h *MlvtController) ProcessFullPipeline(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	sourceLang := c.Query("source_language")
	targetLang := c.Query("target_language")

	// region STT

	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid video ID"})
		return
	}

	if sourceLang == "" || targetLang == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "source_language and target_language are required"})
		return
	}

	video, _, _, err := h.mediaService.GetVideoByID(videoID)
	if err != nil || video == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Video not found"})
		return
	}

	// Create initial Transcription entity
	transcriptionFolder := env.EnvConfig.TranscriptionsFolder
	if transcriptionFolder == "" {
		transcriptionFolder = "transcriptions"
	}
	transcriptionFileName := fmt.Sprintf("transcription_%d.txt", videoID)
	transcription := &entity.Transcription{
		VideoID:   videoID,
		UserID:    video.UserID,
		Folder:    transcriptionFolder,
		FileName:  transcriptionFileName,
		Status:    entity.StatusProcessing,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	transcriptionID, err := h.mediaService.CreateTranscription(transcription, true, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store transcription data"})
		return
	}

	h.quickLogTraffic(entity.ProcessSTTModelAction, video.UserID, transcriptionID)

	// Create output Video entity
	videoFolder := env.EnvConfig.VideosFolder
	if videoFolder == "" {
		videoFolder = "raw_videos"
	}
	outputVideoFileName := fmt.Sprintf("full_pipeline_%d.mp4", videoID)
	outputVideo := &entity.Video{
		OriginalVideoID: videoID,
		UserID:          video.UserID,
		Folder:          videoFolder,
		Image:           video.Image,
		FileName:        outputVideoFileName,
		Status:          entity.StatusProcessing,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	outputVideoID, err := h.mediaService.CreateVideo(outputVideo, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store output video data"})
		return
	}

	h.quickLogTraffic(entity.ProcessLSModelAction, video.UserID, outputVideoID)
	h.quickLogTraffic(entity.ProcessFullPipelineModelAction, video.UserID, outputVideoID)

	// Send notification about Full Pipeline processing start
	h.sendNotification(fmt.Sprintf("🚀 <b>Full Pipeline Processing Started</b>\n\n"+
		"📹 Input Video ID: %d\n"+
		"📹 Output Video ID: %d\n"+
		"👤 User ID: %d\n"+
		"🗣️ From: %s → %s\n\n"+
		"📋 Pipeline Steps:\n"+
		"1️⃣ Speech-to-Text\n"+
		"2️⃣ Text Translation\n"+
		"3️⃣ Text-to-Speech\n"+
		"4️⃣ Lip Sync\n\n"+
		"⏳ Processing started...", videoID, outputVideoID, video.UserID, sourceLang, targetLang))

	// Insert to mongodb
	sttDocument := &entity.Progress{
		UserID:                    video.UserID,
		ProgressType:              entity.ProgressTypeFP,
		OriginalVideoID:           videoID,
		OriginalTranscriptionID:   transcriptionID,
		TranslatedTranscriptionID: 0,
		AudioID:                   0,
		ProgressedVideoID:         outputVideoID,
		Status:                    entity.StatusProcessing,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}

	documentId, err := h.progressService.Create(context.Background(), *sttDocument)
	if err != nil {
		log.Errorf("Failed to insert document ", err)
	}
	log.Infof("Added document STT, Id: ", documentId)

	// Respond immediately to the client
	c.JSON(http.StatusAccepted, response.MessageCreateResponseWithID{
		Message: "Accepted for processing",
		Id:      outputVideoID,
	})

	// Start asynchronous processing
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Warnf("Recovered in goroutine: %v", r)
				h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
				h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)

				// Send failure notification
				h.sendNotification(fmt.Sprintf("❌ <b>Full Pipeline Processing Failed</b>\n\n"+
					"📹 Input Video ID: %d\n"+
					"📹 Output Video ID: %d\n"+
					"👤 User ID: %d\n"+
					"🗣️ Language: %s → %s\n\n"+
					"💥 Error: %v\n\n"+
					"Pipeline failed during processing.", videoID, outputVideoID, video.UserID, sourceLang, targetLang, r))
			}
		}()

		// Step 1: Speech-to-Text
		log.Infof("step 1: speech to text \n")
		h.sendNotification(fmt.Sprintf("1️⃣ <b>Full Pipeline - Step 1: Speech-to-Text</b>\n\n"+
			"📹 Video ID: %d\n"+
			"👤 User ID: %d\n"+
			"🗣️ Language: %s\n\n"+
			"⏳ Converting speech to text...", videoID, video.UserID, sourceLang))
		videoDownloadURL, err := h.mediaService.GeneratePresignedDownloadURLForVideo(videoID)
		if err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate video download URL: %v", err)
			return
		}

		transcriptionUploadURL, err := h.mediaService.GeneratePresignedUploadURLForText(transcriptionFolder, transcriptionFileName, "text/plain")
		if err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate transcription upload URL: %v", err)
			return
		}

		sttPayload := request.STTRequest{
			BaseRequest: request.BaseRequest{
				InputFileName:  video.FileName,
				InputLink:      videoDownloadURL,
				OutputFileName: transcriptionFileName,
				OutputLink:     transcriptionUploadURL,
				Model:          "",
			},
		}

		// Marshal the payload to JSON for curl logging
		sttPayloadBytes, err := json.Marshal(sttPayload)
		if err != nil {
			log.Warnf("Error marshaling Full Pipeline STT payload: %v", err)
		}

		ec2STTURL := fmt.Sprintf("http://%s:%s/stt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)

		// Build and print the curl command for debugging
		curlCmd := fmt.Sprintf(`curl -X POST "%s" -H "Content-Type: application/json" -d '%s'`, ec2STTURL, string(sttPayloadBytes))
		fmt.Println("Full Pipeline STT Curl Command:", curlCmd)

		ec2STTResponse, err := sendRequestToEC2(sttPayload, ec2STTURL, 5*time.Minute)
		if err != nil || ec2STTResponse.Status != "succeeded" {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("EC2 STT processing failed: %v", err)
			return
		}

		h.mediaService.UpdateTranscription(&entity.Transcription{
			ID:        transcriptionID,
			Text:      ec2STTResponse.Result,
			UpdatedAt: time.Now(),
		})

		// Update status to succeeded
		if err := h.mediaService.UpdateTranscriptionStatus(transcriptionID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update transcription status: %v", err)
		}

		// #endregion

		// #region TTT

		// Step 2: Text-to-Text
		log.Infof("step 2: text to text \n")
		h.sendNotification(fmt.Sprintf("2️⃣ <b>Full Pipeline - Step 2: Text Translation</b>\n\n"+
			"📹 Video ID: %d\n"+
			"👤 User ID: %d\n"+
			"🗣️ From: %s → %s\n\n"+
			"⏳ Translating text...", videoID, video.UserID, sourceLang, targetLang))
		translatedFileName := fmt.Sprintf("transcription_%d_%s.txt", transcriptionID, targetLang)
		translatedTranscription := &entity.Transcription{
			VideoID:                 videoID,
			UserID:                  video.UserID,
			OriginalTranscriptionID: transcriptionID,
			Lang:                    targetLang,
			Folder:                  transcriptionFolder,
			FileName:                translatedFileName,
			Status:                  entity.StatusProcessing,
			CreatedAt:               time.Now(),
			UpdatedAt:               time.Now(),
		}

		translatedTranscriptionID, err := h.mediaService.CreateTranscription(translatedTranscription, true, false)
		if err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to create translated transcription: %v", err)
			return
		}

		h.quickLogTraffic(entity.ProcessTTTModelAction, video.UserID, translatedTranscriptionID)

		// Update translated transcription ID to mongodb progress
		h.progressService.UpdateFieldId(context.Background(), documentId, "TranslatedTranscriptionID", translatedTranscriptionID)

		originalDownloadURL, err := h.mediaService.GeneratePresignedDownloadURLForText(transcriptionID)
		if err != nil {
			h.mediaService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate original transcription download URL: %v", err)
			return
		}

		translationUploadURL, err := h.mediaService.GeneratePresignedUploadURLForText(transcriptionFolder, translatedFileName, "text/plain")
		if err != nil {
			h.mediaService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate translation upload URL: %v", err)
			return
		}

		tttPayload := request.TTTRequest{
			BaseRequest: request.BaseRequest{
				InputFileName:  transcriptionFileName,
				InputLink:      originalDownloadURL,
				OutputFileName: translatedFileName,
				OutputLink:     translationUploadURL,
				Model:          "",
			},
			BaseLang: request.BaseLang{
				SourceLang: sourceLang,
				TargetLang: targetLang,
			},
		}

		// Marshal the payload to JSON for curl logging
		tttPayloadBytes, err := json.Marshal(tttPayload)
		if err != nil {
			log.Warnf("Error marshaling Full Pipeline TTT payload: %v", err)
		}

		ec2TTTURL := fmt.Sprintf("http://%s:%s/ttt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)

		// Build and print the curl command for debugging
		tttCurlCmd := fmt.Sprintf(`curl -X POST "%s" -H "Content-Type: application/json" -d '%s'`, ec2TTTURL, string(tttPayloadBytes))
		fmt.Println("Full Pipeline TTT Curl Command:", tttCurlCmd)

		ec2TTTResponse, err := sendRequestToEC2(tttPayload, ec2TTTURL, 5*time.Minute)
		if err != nil || ec2TTTResponse.Status != "succeeded" {
			h.mediaService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("EC2 TTT processing failed: %v", err)
			return
		}

		h.mediaService.UpdateTranscription(&entity.Transcription{
			ID:        translatedTranscriptionID,
			Text:      ec2TTTResponse.Result,
			UpdatedAt: time.Now(),
		})

		// Update status to succeeded
		if err := h.mediaService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update transcription status: %v", err)
		}

		// #endregion

		// region TTS

		// Step 3: Text-to-Speech
		log.Infof("step 3: text to speech \n")
		h.sendNotification(fmt.Sprintf("3️⃣ <b>Full Pipeline - Step 3: Text-to-Speech</b>\n\n"+
			"📹 Video ID: %d\n"+
			"👤 User ID: %d\n"+
			"🗣️ Language: %s\n\n"+
			"⏳ Converting text to speech...", videoID, video.UserID, targetLang))
		audioFolder := env.EnvConfig.AudioFolder
		if audioFolder == "" {
			audioFolder = "audios"
		}
		audioFileName := fmt.Sprintf("audio_%d.mp3", translatedTranscriptionID)
		audio := &entity.Audio{
			TranscriptionID: translatedTranscriptionID,
			VideoID:         videoID,
			UserID:          video.UserID,
			Lang:            targetLang,
			Folder:          audioFolder,
			FileName:        audioFileName,
			Status:          entity.StatusProcessing,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		audioID, err := h.mediaService.CreateAudio(audio, true)
		if err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to create audio: %v", err)
			return
		}

		h.quickLogTraffic(entity.ProcessTTSModelAction, video.UserID, audioID)

		// update audio ID to mongodb progress collection
		h.progressService.UpdateFieldId(context.Background(), documentId, "AudioID", audioID)

		transcriptionDownloadURL, err := h.mediaService.GeneratePresignedDownloadURLForText(translatedTranscriptionID)
		if err != nil {
			h.mediaService.UpdateAudioStatus(audioID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate transcription download URL: %v", err)
			return
		}

		audioUploadURL, err := h.mediaService.GeneratePresignedUploadURL(audioFolder, audioFileName, "audio/mpeg")
		if err != nil {
			h.mediaService.UpdateAudioStatus(audioID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate audio upload URL: %v", err)
			return
		}

		ttsPayload := request.TTSRequest{
			BaseRequest: request.BaseRequest{
				InputFileName:  translatedFileName,
				InputLink:      transcriptionDownloadURL,
				OutputFileName: audioFileName,
				OutputLink:     audioUploadURL,
				Model:          "",
			},
			InputAudioFileName: video.FileName,
			InputAudioLink:     videoDownloadURL,
			Lang:               targetLang,
		}

		// Marshal the payload to JSON for curl logging
		ttsPayloadBytes, err := json.Marshal(ttsPayload)
		if err != nil {
			log.Warnf("Error marshaling Full Pipeline TTS payload: %v", err)
		}

		ec2TTSURL := fmt.Sprintf("http://%s:%s/tts", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)

		// Build and print the curl command for debugging
		ttsCurlCmd := fmt.Sprintf(`curl -X POST "%s" -H "Content-Type: application/json" -d '%s'`, ec2TTSURL, string(ttsPayloadBytes))
		fmt.Println("Full Pipeline TTS Curl Command:", ttsCurlCmd)

		ec2TTSResponse, err := sendRequestToEC2(ttsPayload, ec2TTSURL, 5*time.Minute)
		if err != nil || ec2TTSResponse.Status != "succeeded" {
			h.mediaService.UpdateAudioStatus(audioID, entity.StatusFailed)
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("EC2 TTS processing failed: %v", err)
			return
		}

		h.mediaService.UpdateAudio(&entity.Audio{
			ID:        audioID,
			UpdatedAt: time.Now(),
		})

		// Update status to succeeded
		if err := h.mediaService.UpdateAudioStatus(audioID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update audio status: %v", err)
		}

		// #endregion

		// #region LS

		// Step 4: Lip Sync
		log.Infof("step 4: lipsync \n")
		h.sendNotification(fmt.Sprintf("4️⃣ <b>Full Pipeline - Step 4: Lip Sync</b>\n\n"+
			"📹 Video ID: %d\n"+
			"👤 User ID: %d\n"+
			"🎧 Audio ID: %d\n\n"+
			"⏳ Synchronizing lip movements...", videoID, video.UserID, audioID))
		outputVideo.AudioID = audioID
		outputVideo.ID = outputVideoID
		log.Warnf("error: %v \n", outputVideo)
		if err := h.mediaService.UpdateVideo(outputVideo); err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to update output video: %v", err)
			return
		}

		videoDownloadURL, err = h.mediaService.GeneratePresignedDownloadURLForVideo(videoID)
		if err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate video download URL: %v", err)
			return
		}

		audioDownloadURL, err := h.mediaService.GeneratePresignedDownloadURL(audioID)
		if err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate audio download URL: %v", err)
			return
		}

		outputVideoUploadURL, err := h.mediaService.GeneratePresignedUploadURLForVideo(videoFolder, outputVideoFileName, "video/mp4")
		if err != nil {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate output video upload URL: %v", err)
			return
		}

		lsPayload := request.LSRequest{
			InputVideoFileName:  video.FileName,
			InputVideoLink:      videoDownloadURL,
			InputAudioFileName:  audioFileName,
			InputAudioLink:      audioDownloadURL,
			OutputVideoFileName: outputVideoFileName,
			OutputVideoLink:     outputVideoUploadURL,
			Model:               "",
		}

		// Marshal the payload to JSON
		payloadBytes, err := json.Marshal(lsPayload)
		if err != nil {
			log.Warnf("Error marshaling payload: %v", err)
		}

		// Construct the EC2 URL (adjust your EC2 IP and Port as needed)
		ec2URL := fmt.Sprintf("http://%s:%s/ls", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)

		// Build the curl command string
		curlCmd = fmt.Sprintf(`curl -X POST "%s" -H "Content-Type: application/json" -d '%s'`, ec2URL, string(payloadBytes))

		// Print the curl command to the console
		fmt.Println(curlCmd)

		ec2LSURL := fmt.Sprintf("http://%s:%s/ls", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		ec2LSResponse, err := sendRequestToEC2(lsPayload, ec2LSURL, 25*time.Minute)
		if err != nil || ec2LSResponse.Status != "succeeded" {
			h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("EC2 Lip Sync processing failed: %v", err)
			return
		}

		h.mediaService.UpdateVideo(&entity.Video{
			ID:        outputVideoID,
			UpdatedAt: time.Now(),
		})

		// Update status to succeeded
		if err := h.mediaService.UpdateVideoStatus(outputVideoID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update video status: %v", err)
		}

		// #endregion

		h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusSucceeded)

		// Send final success notification
		h.sendNotification(fmt.Sprintf("🎉 <b>Full Pipeline Processing Completed Successfully!</b>\n\n"+
			"📹 Input Video ID: %d\n"+
			"📹 Output Video ID: %d\n"+
			"👤 User ID: %d\n"+
			"🗣️ Language: %s → %s\n"+
			"📄 Output File: %s\n\n"+
			"✅ All steps completed:\n"+
			"✓ Speech-to-Text\n"+
			"✓ Text Translation\n"+
			"✓ Text-to-Speech\n"+
			"✓ Lip Sync\n\n"+
			"🚀 Your multilingual video is ready!", videoID, outputVideoID, video.UserID, sourceLang, targetLang, outputVideoFileName))

		log.Infof("Finish: fullpipeline \n")
	}()
}
