package progress_handler

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/media_service"
	"mlvt/internal/service/progress_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgressController struct {
	progressService progress_service.ProgressService
	mediaService    media_service.MediaService
}

func NewProgressService(
	progressService progress_service.ProgressService,
	mediaService media_service.MediaService,
) *ProgressController {
	return &ProgressController{
		progressService: progressService,
		mediaService:    mediaService,
	}
}

// GetUserProgress godoc
// @Summary Get user progress
// @Description Retrieves the progress data for a specified user by user ID
// @Tags Progress
// @Accept  json
// @Produce  json
// @Param   user_id path     uint64 true "User ID"
// @Success 200 {object} map[string]interface{} "A JSON object containing progress data. Key: 'progresses'"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /progress/{user_id} [get]
func (h *ProgressController) GetUserProgress(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userId, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid userId"})
		return
	}

	listProgress := []entity.ProgressType{
		entity.ProgressTypeFP,
		entity.ProgressTypeTTS,
		entity.ProgressTypeSTT,
		entity.ProgressTypeLS,
		entity.ProgressTypeTTT,
	}
	listStatus := []entity.StatusEntity{
		entity.StatusFailed,
		entity.StatusProcessing,
		entity.StatusRaw,
		entity.StatusSucceeded,
	}
	progresses, _, err := h.progressService.GetProgressByUserID(context.Background(), userId, 0, 15, "", listProgress, listStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to get user progress"})
		log.Errorf("failed to get user progress, err: ", err)
		return
	}

	progressWithThumbnail, err := h.progressService.GetProgressThumbnails(progresses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to add thumbnail to progress entity"})
		log.Errorf("failed to add thumbnail to progress entity")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"progresses": progressWithThumbnail,
	})
}

// UpdateProgressTitle godoc
// @Summary Update progress title
// @Description Updates the title of a specific progress by progress ID
// @Tags Progress
// @Accept  json
// @Produce  json
// @Param   progress_id path     string true "Progress ID"
// @Param   request body map[string]string true "Request body containing title"
// @Success 200 {object} map[string]interface{} "Title updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid progress ID or missing title"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /progress/{progress_id}/title [put]
func (h *ProgressController) UpdateProgressTitle(c *gin.Context) {
	progressID := c.Param("progress_id")
	id, err := primitive.ObjectIDFromHex(progressID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid progress id"})
		return
	}

	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	title, ok := body["title"]
	if !ok || title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing title"})
		return
	}

	if err := h.progressService.UpdateTitle(context.Background(), id, title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update title"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Title updated"})
}

// DeleteProgress godoc
// @Summary Delete progress
// @Description Deletes a specific progress and its associated media files by progress ID
// @Tags Progress
// @Accept  json
// @Produce  json
// @Param   progress_id path     string true "Progress ID"
// @Success 200 {object} map[string]interface{} "Progress deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid progress ID or progress not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /progress/{progress_id} [delete]
func (h *ProgressController) DeleteProgress(c *gin.Context) {
	progressID := c.Param("progress_id")
	id, err := primitive.ObjectIDFromHex(progressID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid progress id"})
		return
	}

	progressInfo, err := h.progressService.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Progress not found"})
		return
	}

	if err := h.progressService.DeleteProgress(context.Background(), id); err != nil {
		log.Errorf("failed to delete progress: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete progress"})
		return
	}

	if progressInfo.ProgressType == entity.ProgressTypeTTS {
		// delete audio
		err := h.mediaService.DeleteAudio(progressInfo.AudioID)
		if err != nil {
			log.Errorf("failed to delete result audio: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete result audio"})
			return
		}
	} else {
		if progressInfo.ProgressedVideoID != 0 {
			err := h.mediaService.DeleteVideo(progressInfo.ProgressedVideoID)
			if err != nil {
				log.Errorf("failed to delete result video: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete result video"})
				return
			}
		}
		if progressInfo.TranslatedTranscriptionID != 0 {
			err := h.mediaService.DeleteTranscription(progressInfo.TranslatedTranscriptionID)
			if err != nil {
				log.Errorf("failed to delete result transcription: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete result transcription"})
				return
			}
		}
		if progressInfo.ProgressType == entity.ProgressTypeFP && progressInfo.OriginalTranscriptionID != 0 {
			err := h.mediaService.DeleteTranscription(progressInfo.TranslatedTranscriptionID)
			if err != nil {
				log.Errorf("failed to delete extracted transcription: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete extracted transcription"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progress deleted"})
}
