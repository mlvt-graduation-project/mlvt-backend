package handler

import (
	"context"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PingController struct {
	Service service.PingService
}

func NewPingController(service service.PingService) *PingController {
	return &PingController{
		Service: service,
	}
}

// Handler functions

func (h *PingController) PingSpeechToText(c *gin.Context) {
	h.ping(c, h.Service.PingSpeechToText)
}

func (h *PingController) PingTextToText(c *gin.Context) {
	h.ping(c, h.Service.PingTextToText)
}

func (h *PingController) PingTextToSpeech(c *gin.Context) {
	h.ping(c, h.Service.PingTextToSpeech)
}

func (h *PingController) PingVoiceCloning(c *gin.Context) {
	h.ping(c, h.Service.PingVoiceCloning)
}

func (h *PingController) PingLipSync(c *gin.Context) {
	h.ping(c, h.Service.PingLipSync)
}

func (h *PingController) PingFullPipeline(c *gin.Context) {
	h.ping(c, h.Service.PingFullPipeline)
}

// Utility method to reduce duplication
type pingFunc func(ctx context.Context, id uint64) (*response.PingStatusResponse, error)

func (h *PingController) ping(c *gin.Context, pf pingFunc) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()

	status, err := pf(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
