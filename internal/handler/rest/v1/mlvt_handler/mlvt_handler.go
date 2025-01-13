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
	"mlvt/internal/service/audio_service"
	"mlvt/internal/service/progress_service"
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
	progressService      progress_service.ProgressService
}

func NewMlvtController(
	audioService audio_service.AudioService,
	transcriptionService transcription_service.TranscriptionService,
	videoService video_service.VideoService,
	progressService progress_service.ProgressService,
) *MlvtController {
	return &MlvtController{
		audioService:         audioService,
		transcriptionService: transcriptionService,
		videoService:         videoService,
		progressService:      progressService,
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

	video, _, _, err := h.videoService.GetVideoByID(videoID)
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
		FileName:  transcriptionFileName,
		Status:    entity.StatusProcessing,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	transcriptionID, err := h.transcriptionService.CreateTranscription(transcription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store transcription data"})
		return
	}

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
				h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed)
				h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			}
		}()

		videoDownloadURL, err := h.videoService.GeneratePresignedDownloadURLForVideo(videoID)
		if err != nil {
			h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate video download URL: %v", err)
			return
		}

		fileType := "text/plain"
		transcriptionUploadURL, err := h.transcriptionService.GeneratePresignedUploadURL(folder, transcriptionFileName, fileType)
		if err != nil {
			h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed)
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

		ec2ServerURL := fmt.Sprintf("http://%s:%s/stt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		ec2Response, err := sendRequestToEC2(requestPayload, ec2ServerURL, 5*time.Minute)
		log.Infof("ec2 response: %v\n\n", ec2Response)
		if err != nil || ec2Response.Status != "succeeded" {
			h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("EC2 processing failed: %v", err)
			return
		}

		updateTranscription := &entity.Transcription{
			ID:        transcriptionID,
			Text:      ec2Response.Result,
			UpdatedAt: time.Now(),
		}

		if err := h.transcriptionService.UpdateTranscription(updateTranscription); err != nil {
			log.Errorf("Failed to update transcription data: %v", err)
		}

		// Update status to succeeded
		if err := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update transcription status: %v", err)
		}

		h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusSucceeded)
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

	originalTranscription, _, err := h.transcriptionService.GetTranscriptionByID(transcriptionID)
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

	translatedTranscriptionID, err := h.transcriptionService.CreateTranscription(newTranscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store translated transcription data"})
		return
	}

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
				h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
				h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			}
		}()

		originalDownloadURL, err := h.transcriptionService.GeneratePresignedDownloadURL(transcriptionID)
		if err != nil {
			h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("Failed to generate original transcription download URL: %v", err)
			return
		}

		fileType := "text/plain"
		translationUploadURL, err := h.transcriptionService.GeneratePresignedUploadURL(folder, translatedFileName, fileType)
		if err != nil {
			h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
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

		ec2ServerURL := fmt.Sprintf("http://%s:%s/ttt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		ec2Response, err := sendRequestToEC2(requestPayload, ec2ServerURL, 5*time.Minute)
		if err != nil || ec2Response.Status != "succeeded" {
			h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusFailed)
			log.Errorf("EC2 processing failed: %v", err)
			return
		}

		updateTranscription := &entity.Transcription{
			ID:        translatedTranscriptionID,
			Text:      ec2Response.Result,
			UpdatedAt: time.Now(),
		}

		if err := h.transcriptionService.UpdateTranscription(updateTranscription); err != nil {
			log.Errorf("Failed to update transcription data: %v", err)
		}

		// Update status to succeeded
		if err := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update transcription status: %v", err)
		}

		h.progressService.UpdateStatus(context.Background(), documentId, entity.StatusSucceeded)
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

	transcription, _, err := h.transcriptionService.GetTranscriptionByID(transcriptionID)
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

	audioID, err := h.audioService.CreateAudio(audio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store audio data"})
		return
	}

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
				h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed)
			}
		}()

		transcriptionDownloadURL, err := h.transcriptionService.GeneratePresignedDownloadURL(transcriptionID)
		if err != nil {
			h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed)
			log.Errorf("Failed to generate transcription download URL: %v", err)
			return
		}

		fileType := "audio/mpeg"
		audioUploadURL, err := h.audioService.GeneratePresignedUploadURL(folder, audioFileName, fileType)
		if err != nil {
			h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed)
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
			Lang: transcription.Lang,
		}

		ec2ServerURL := fmt.Sprintf("http://%s:%s/tts", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		ec2Response, err := sendRequestToEC2(requestPayload, ec2ServerURL, 5*time.Minute)
		if err != nil || ec2Response.Status != "succeeded" {
			h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed)
			log.Errorf("EC2 processing failed: %v", err)
			return
		}

		updateAudio := &entity.Audio{
			ID:        audioID,
			UpdatedAt: time.Now(),
		}

		if err := h.audioService.UpdateAudio(updateAudio); err != nil {
			log.Errorf("Failed to update audio data: %v", err)
		}
		// Update status to succeeded
		if err := h.audioService.UpdateAudioStatus(audioID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update audio status: %v", err)
		}
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

	video, _, _, err := h.videoService.GetVideoByID(videoID)
	if err != nil || video == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Video not found"})
		return
	}

	audio, _, err := h.audioService.GetAudioByID(audioID)
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
		FileName:        outputVideoFileName,
		Status:          entity.StatusProcessing,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	outputVideoID, err := h.videoService.CreateVideo(outputVideo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store output video data"})
		return
	}

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
				h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			}
		}()

		videoDownloadURL, err := h.videoService.GeneratePresignedDownloadURLForVideo(videoID)
		if err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("Failed to generate video download URL: %v", err)
			return
		}

		audioDownloadURL, err := h.audioService.GeneratePresignedDownloadURL(audioID)
		if err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("Failed to generate audio download URL: %v", err)
			return
		}

		fileType := "video/mp4"
		outputVideoUploadURL, err := h.videoService.GeneratePresignedUploadURLForVideo(folder, outputVideoFileName, fileType)
		if err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
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

		ec2ServerURL := fmt.Sprintf("http://%s:%s/ls", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		ec2Response, err := sendRequestToEC2(requestPayload, ec2ServerURL, 15*time.Minute)
		if err != nil || ec2Response.Status != "succeeded" {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("EC2 processing failed: %v", err)
			return
		}

		updateVideo := &entity.Video{
			ID:        outputVideoID,
			UpdatedAt: time.Now(),
		}

		if err := h.videoService.UpdateVideo(updateVideo); err != nil {
			log.Errorf("Failed to update video data: %v", err)
		}

		// Update status to succeeded
		if err := h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update video status: %v", err)
		}
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

	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid video ID"})
		return
	}

	if sourceLang == "" || targetLang == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "source_language and target_language are required"})
		return
	}

	video, _, _, err := h.videoService.GetVideoByID(videoID)
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

	transcriptionID, err := h.transcriptionService.CreateTranscription(transcription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store transcription data"})
		return
	}

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

	outputVideoID, err := h.videoService.CreateVideo(outputVideo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to store output video data"})
		return
	}

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
				h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			}
		}()

		// Step 1: Speech-to-Text
		log.Infof("step 1: speech to text \n")
		videoDownloadURL, err := h.videoService.GeneratePresignedDownloadURLForVideo(videoID)
		if err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("Failed to generate video download URL: %v", err)
			return
		}

		transcriptionUploadURL, err := h.transcriptionService.GeneratePresignedUploadURL(transcriptionFolder, transcriptionFileName, "text/plain")
		if err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
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

		ec2STTURL := fmt.Sprintf("http://%s:%s/stt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		ec2STTResponse, err := sendRequestToEC2(sttPayload, ec2STTURL, 5*time.Minute)
		if err != nil || ec2STTResponse.Status != "succeeded" {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("EC2 STT processing failed: %v", err)
			return
		}

		h.transcriptionService.UpdateTranscription(&entity.Transcription{
			ID:        transcriptionID,
			Text:      ec2STTResponse.Result,
			UpdatedAt: time.Now(),
		})

		// Update status to succeeded
		if err := h.transcriptionService.UpdateTranscriptionStatus(transcriptionID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update transcription status: %v", err)
		}

		// Step 2: Text-to-Text
		log.Infof("step 2: text to text \n")
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

		translatedTranscriptionID, err := h.transcriptionService.CreateTranscription(translatedTranscription)
		if err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("Failed to create translated transcription: %v", err)
			return
		}

		originalDownloadURL, err := h.transcriptionService.GeneratePresignedDownloadURL(transcriptionID)
		if err != nil {
			h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			log.Errorf("Failed to generate original transcription download URL: %v", err)
			return
		}

		translationUploadURL, err := h.transcriptionService.GeneratePresignedUploadURL(transcriptionFolder, translatedFileName, "text/plain")
		if err != nil {
			h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
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

		ec2TTTURL := fmt.Sprintf("http://%s:%s/ttt", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		ec2TTTResponse, err := sendRequestToEC2(tttPayload, ec2TTTURL, 5*time.Minute)
		if err != nil || ec2TTTResponse.Status != "succeeded" {
			h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusFailed)
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("EC2 TTT processing failed: %v", err)
			return
		}

		h.transcriptionService.UpdateTranscription(&entity.Transcription{
			ID:        translatedTranscriptionID,
			Text:      ec2TTTResponse.Result,
			UpdatedAt: time.Now(),
		})

		// Update status to succeeded
		if err := h.transcriptionService.UpdateTranscriptionStatus(translatedTranscriptionID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update transcription status: %v", err)
		}

		// Step 3: Text-to-Speech
		log.Infof("step 3: text to speech \n")
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

		audioID, err := h.audioService.CreateAudio(audio)
		if err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("Failed to create audio: %v", err)
			return
		}

		transcriptionDownloadURL, err := h.transcriptionService.GeneratePresignedDownloadURL(translatedTranscriptionID)
		if err != nil {
			h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed)
			log.Errorf("Failed to generate transcription download URL: %v", err)
			return
		}

		audioUploadURL, err := h.audioService.GeneratePresignedUploadURL(audioFolder, audioFileName, "audio/mpeg")
		if err != nil {
			h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed)
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
			Lang: targetLang,
		}

		ec2TTSURL := fmt.Sprintf("http://%s:%s/tts", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		ec2TTSResponse, err := sendRequestToEC2(ttsPayload, ec2TTSURL, 5*time.Minute)
		if err != nil || ec2TTSResponse.Status != "succeeded" {
			h.audioService.UpdateAudioStatus(audioID, entity.StatusFailed)
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("EC2 TTS processing failed: %v", err)
			return
		}

		h.audioService.UpdateAudio(&entity.Audio{
			ID:        audioID,
			UpdatedAt: time.Now(),
		})

		// Update status to succeeded
		if err := h.audioService.UpdateAudioStatus(audioID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update audio status: %v", err)
		}

		// Step 4: Lip Sync
		log.Infof("step 4: lipsync \n")
		outputVideo.AudioID = audioID
		outputVideo.ID = outputVideoID
		log.Warnf("error: %v \n", outputVideo)
		if err := h.videoService.UpdateVideo(outputVideo); err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("Failed to update output video: %v", err)
			return
		}

		videoDownloadURL, err = h.videoService.GeneratePresignedDownloadURLForVideo(videoID)
		if err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("Failed to generate video download URL: %v", err)
			return
		}

		audioDownloadURL, err := h.audioService.GeneratePresignedDownloadURL(audioID)
		if err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("Failed to generate audio download URL: %v", err)
			return
		}

		outputVideoUploadURL, err := h.videoService.GeneratePresignedUploadURLForVideo(videoFolder, outputVideoFileName, "video/mp4")
		if err != nil {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
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

		ec2LSURL := fmt.Sprintf("http://%s:%s/ls", env.EnvConfig.Ec2IPAddress, env.EnvConfig.Ec2Port)
		ec2LSResponse, err := sendRequestToEC2(lsPayload, ec2LSURL, 15*time.Minute)
		if err != nil || ec2LSResponse.Status != "succeeded" {
			h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusFailed)
			log.Errorf("EC2 Lip Sync processing failed: %v", err)
			return
		}

		h.videoService.UpdateVideo(&entity.Video{
			ID:        outputVideoID,
			UpdatedAt: time.Now(),
		})

		// Update status to succeeded
		if err := h.videoService.UpdateVideoStatus(outputVideoID, entity.StatusSucceeded); err != nil {
			log.Errorf("Failed to update video status: %v", err)
		}

		log.Infof("Finish: fullpipeline \n")
	}()
}
