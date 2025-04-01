package admin_handler

import (
	"context"
	"fmt"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/request"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *AdminController) GetServerConfig(c *gin.Context) {
	ctx := context.Background()

	config, err := h.adminService.GetServerConfig(ctx)
	if err != nil {
		log.Errorf("error when getting server config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get server config"})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *AdminController) UpdateServerConfig(c *gin.Context) {
	ctx := context.Background()

	adminIdStr := c.Param("adminID")
	adminId, err := strconv.ParseUint(adminIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	var req request.UpdateServerConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	err = h.adminService.UpdateServerConfig(ctx, adminId, req.ModelType, req.ModelName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to update server config: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "server config updated"})
}
