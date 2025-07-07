package process_handler

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/request"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/media_service"
	"mlvt/internal/service/progress_service"
	"mlvt/internal/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProcessController struct {
	processServivce progress_service.ProgressService
	mediaService media_service.MediaService
}

func NewProcessService (processService progress_service.ProgressService, mediaService media_service.MediaService) *ProcessController{
	return &ProcessController{
		processServivce: processService,
		mediaService: mediaService,
	}
}

func (h *ProcessController) GetAllProcess (c *gin.Context) {
	userIDStr := c.Param("user_id")
	userId, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid userId"})
		return
	}

	var request request.ProcessRequest
	var result []response.ProcessResponse
	err = c.ShouldBindBodyWithJSON(&request) 
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid request body"})
	}

	// Get all project process
	if len(request.ProjectType) > 0 {
		// Get project from mongodb
		progresses, err := h.processServivce.GetProgressByUserID(context.Background(), userId, request.Offset, request.Limit, request.SearchKey, request.ProjectType, request.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to get user progress"})
			log.Errorf("failed to get user progress, err: ", err)
			return
		}

		// Add thumbnail to projects
		progressWithThumbnail, err := h.processServivce.GetProgressThumbnails(progresses)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to add thumbnail to progress entity"})
			log.Errorf("failed to add thumbnail to progress entity")
			return
		}
		listProcess := utility.ProgressResponseListToProcessResponseList(progressWithThumbnail)
		result = append(result, listProcess...)
	}
	
	// Get Audio process
	if utility.Contains(request.MediaType, entity.MediaTypeAudio) {
		audio, err := h.mediaService.ListAudiosByUserIDAdvance(userId, request.SearchKey, request.Limit, request.Offset, request.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to get list audio"})
			log.Errorf("GetAllProcess: Failed to get list audio by user id: %v", err)
			return
		}
		audioProcesses := utility.AudioListToProcessResponseList(audio)
		result = append(result, audioProcesses...)
	} 
	
	// Get Video process
	if utility.Contains(request.MediaType, entity.MediaTypeVideo) {
		video, err := h.mediaService.ListVideosByUserIDAdvance(userId, request.SearchKey, request.Limit, request.Offset, request.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to get list video"})
			log.Errorf("GetAllProcess: Failed to get list video by user id: %v", err)
			return
		}
		videoProcess := utility.VideoResponseListToProcessResponseList(video)
		result = append(result, videoProcess...)
	}
	
	// Get text process
	if utility.Contains(request.MediaType, entity.MediaTypeText) {
		text, err := h.mediaService.ListTranscriptionsByUserIDAdvance(userId, request.SearchKey, request.Limit, request.Offset, request.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to get transcription video"})
			log.Errorf("GetAllProcess: Failed to get list transcription by user id: %v", err)
			return
		}
		textProcess := utility.TranscriptionListToProcessResponseList(text)
		result = append(result, textProcess...)
	}
	
	response := utility.SortProcessResponseByCreateDate(result, request.Limit)
	c.JSON(http.StatusOK, response)

	return 
}