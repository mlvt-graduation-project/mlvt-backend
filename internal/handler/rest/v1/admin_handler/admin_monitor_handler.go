package admin_handler

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *AdminController) GetMonitorDataType(c *gin.Context) {
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

func (h *AdminController) GetMonitorPipeline(c *gin.Context) {
	ctx := context.Background()

	// parse adminID from URL param
	adminIDStr := c.Param("adminID")
	adminID, err := strconv.ParseUint(adminIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	// call service
	pipeline, err := h.adminService.GetMonitorPipeline(ctx, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get pipeline metrics: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, pipeline)
}

type MonitorTrafficRequest struct {
	CurrentDay string `json:"current_day"`
	Type       string `json:"type"`
}

func (h *AdminController) GetMonitorTraffic(c *gin.Context) {
	ctx := context.Background()

	adminIDStr := c.Param("adminID")
	adminID, err := strconv.ParseUint(adminIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	var req MonitorTrafficRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	var baseTime time.Time
	if req.CurrentDay != "" {
		baseTime, err = time.Parse("2006-01-02", req.CurrentDay)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
	} else {
		baseTime = time.Now()
	}

	// convert string to typed constant
	var periodType entity.TimePeriodType
	switch req.Type {
	case string(entity.TimePeriodDay):
		periodType = entity.TimePeriodDay
	case string(entity.TimePeriodWeek):
		periodType = entity.TimePeriodWeek
	case string(entity.TimePeriodYear):
		periodType = entity.TimePeriodYear
	default:
		// invalid
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time period"})
		return
	}

	trafficData, err := h.adminService.GetMonitorTraffic(ctx, adminID, periodType, baseTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trafficData)
}
