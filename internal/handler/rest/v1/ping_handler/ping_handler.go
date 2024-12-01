package ping_handler

import (
	"context"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/ping_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PingController handles all ping-related API requests.
type PingController struct {
	Service ping_service.PingService
}

// NewPingController creates a new instance of PingController.
func NewPingController(service ping_service.PingService) *PingController {
	return &PingController{
		Service: service,
	}
}

// pingFunc defines the function signature for ping operations.
type pingFunc func(ctx context.Context, id uint64) (*response.PingStatusResponse, error)

/**
 * @Summary Ping Speech-to-Text Status
 * @Description Retrieves the status of a speech-to-text transcription task by ID.
 * @Tags Ping
 * @Accept json
 * @Produce json
 * @Param id path uint64 true "Task ID"
 * @Success 200 {object} response.PingStatusResponse
 * @Failure 400 {object} response.ErrorResponse "Invalid ID"
 * @Failure 500 {object} response.ErrorResponse "Internal Server Error"
 * @Router /ping/speech-to-text/{id} [get]
 */
func (h *PingController) PingSpeechToText(c *gin.Context) {
	h.ping(c, h.Service.PingSpeechToText)
}

/**
 * @Summary Ping Text-to-Text Status
 * @Description Retrieves the status of a text-to-text transcription task by ID.
 * @Tags Ping
 * @Accept json
 * @Produce json
 * @Param id path uint64 true "Task ID"
 * @Success 200 {object} response.PingStatusResponse
 * @Failure 400 {object} response.ErrorResponse "Invalid ID"
 * @Failure 500 {object} response.ErrorResponse "Internal Server Error"
 * @Router /ping/text-to-text/{id} [get]
 */
func (h *PingController) PingTextToText(c *gin.Context) {
	h.ping(c, h.Service.PingTextToText)
}

/**
 * @Summary Ping Text-to-Speech Status
 * @Description Retrieves the status of a text-to-speech audio task by ID.
 * @Tags Ping
 * @Accept json
 * @Produce json
 * @Param id path uint64 true "Task ID"
 * @Success 200 {object} response.PingStatusResponse
 * @Failure 400 {object} response.ErrorResponse "Invalid ID"
 * @Failure 500 {object} response.ErrorResponse "Internal Server Error"
 * @Router /ping/text-to-speech/{id} [get]
 */
func (h *PingController) PingTextToSpeech(c *gin.Context) {
	h.ping(c, h.Service.PingTextToSpeech)
}

/**
 * @Summary Ping Voice Cloning Status
 * @Description Retrieves the status of a voice cloning audio task by ID.
 * @Tags Ping
 * @Accept json
 * @Produce json
 * @Param id path uint64 true "Task ID"
 * @Success 200 {object} response.PingStatusResponse
 * @Failure 400 {object} response.ErrorResponse "Invalid ID"
 * @Failure 500 {object} response.ErrorResponse "Internal Server Error"
 * @Router /ping/voice-cloning/{id} [get]
 */
func (h *PingController) PingVoiceCloning(c *gin.Context) {
	h.ping(c, h.Service.PingVoiceCloning)
}

/**
 * @Summary Ping Lip Sync Status
 * @Description Retrieves the status of a lip-sync video task by ID.
 * @Tags Ping
 * @Accept json
 * @Produce json
 * @Param id path uint64 true "Task ID"
 * @Success 200 {object} response.PingStatusResponse
 * @Failure 400 {object} response.ErrorResponse "Invalid ID"
 * @Failure 500 {object} response.ErrorResponse "Internal Server Error"
 * @Router /ping/lipsync/{id} [get]
 */
func (h *PingController) PingLipSync(c *gin.Context) {
	h.ping(c, h.Service.PingLipSync)
}

/**
 * @Summary Ping Full Pipeline Status
 * @Description Retrieves the status of a full pipeline video task by ID.
 * @Tags Ping
 * @Accept json
 * @Produce json
 * @Param id path uint64 true "Task ID"
 * @Success 200 {object} response.PingStatusResponse
 * @Failure 400 {object} response.ErrorResponse "Invalid ID"
 * @Failure 500 {object} response.ErrorResponse "Internal Server Error"
 * @Router /ping/full-pipeline/{id} [get]
 */
func (h *PingController) PingFullPipeline(c *gin.Context) {
	h.ping(c, h.Service.PingFullPipeline)
}

// ping is a utility method to handle common logic for ping endpoints.
func (h *PingController) ping(c *gin.Context, pf pingFunc) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid ID"})
		return
	}

	ctx := c.Request.Context()

	status, err := pf(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
