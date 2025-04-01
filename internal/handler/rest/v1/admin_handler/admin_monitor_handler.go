package admin_handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *AdminController) GetMonitorDataTypeD(c *gin.Context) {
	ctx := context.Background()

	// Parse the "adminID" (or requestor ID) from URL params
	adminIDStr := c.Param("adminID")
	adminID, err := strconv.ParseUint(adminIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	// Retrieve data via service
	monitorData, err := h.adminService.GetMonitorDataType(ctx, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get monitor data by user ID: %v", err)})
		return
	}

	c.JSON(http.StatusOK, monitorData)
}
