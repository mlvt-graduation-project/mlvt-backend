package process_handler

import (
	"context"
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
	mediaService    media_service.MediaService
}

func NewProcessService(processService progress_service.ProgressService, mediaService media_service.MediaService) *ProcessController {
	return &ProcessController{
		processServivce: processService,
		mediaService:    mediaService,
	}
}

func (h *ProcessController) GetAllProcess(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userId, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid userId"})
		return
	}

	var request request.ProcessRequest
	var result response.ProcessResponse
	err = c.ShouldBindBodyWithJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid request body"})
	}

	// Get all project process
	if len(request.ProjectType) > 0 {
		// Get project from mongodb
		progresses, totalCount, err := h.processServivce.GetProgressByUserID(context.Background(), userId, request.Offset, request.Limit, request.SearchKey, request.ProjectType, request.Status)
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
		result.ProcessList = append(result.ProcessList, listProcess...)
		result.TotalCount = totalCount
	} else if len(request.MediaType) > 0 {
		result, err = h.mediaService.GetAllMedia(userId, request.SearchKey, request.Limit, request.Offset, request.Status, request.MediaType)
		if err != nil {
			log.Errorf("failed to get media. Error: %v", err)
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to get media"})
			return
		}
	}

	c.JSON(http.StatusOK, result)

	return
}
