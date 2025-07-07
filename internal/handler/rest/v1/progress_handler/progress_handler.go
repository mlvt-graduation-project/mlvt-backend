package progress_handler

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/progress_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProgressController struct {
	progressService progress_service.ProgressService
}

func NewProgressService(
	progressService progress_service.ProgressService,
) *ProgressController {
	return &ProgressController{
		progressService: progressService,
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

	listProgress := []entity.ProgressType {
		entity.ProgressTypeFP,
		entity.ProgressTypeTTS,
		entity.ProgressTypeSTT,
		entity.ProgressTypeLS,
		entity.ProgressTypeTTT,
	}
	listStatus := []entity.StatusEntity {
		entity.StatusFailed,
		entity.StatusProcessing,
		entity.StatusRaw,
		entity.StatusSucceeded,
	}
	progresses, err := h.progressService.GetProgressByUserID(context.Background(), userId, 0, 15, "", listProgress, listStatus)
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
