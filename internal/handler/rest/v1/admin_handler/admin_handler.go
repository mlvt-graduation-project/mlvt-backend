package admin_handler

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/service/admin_service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminController struct {
	adminService admin_service.AdminService
}

func NewAdminController(adminService admin_service.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

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

type updateServerConfigRequest struct {
	ModelType string `json:"model_type" binding:"required"`
	ModelName string `json:"model_name" binding:"required"`
}

func (h *AdminController) UpdateServerConfig(c *gin.Context) {
	ctx := context.Background()

	adminIdStr := c.Param("adminID")
	adminId, err := strconv.ParseUint(adminIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	var req updateServerConfigRequest
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

func (h *AdminController) GetModelList(c *gin.Context) {
	ctx := context.Background()

	adminIdStr := c.Param("adminID")
	adminId, err := strconv.ParseUint(adminIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	qo := mongodb.QueryOptions{}

	models, err := h.adminService.GetModelList(ctx, adminId, qo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models)
}

func (h *AdminController) AddModelOption(c *gin.Context) {
	ctx := context.Background()

	adminIdStr := c.Param("adminID")
	adminId, err := strconv.ParseUint(adminIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	var modelOption entity.ModelOption
	if err := c.ShouldBindJSON(&modelOption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	modelOption.ID = primitive.NilObjectID
	modelOption.UpdatedAt = time.Now()

	insertedId, err := h.adminService.AddModelOption(ctx, adminId, modelOption)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "model option added successfully",
		"id":      insertedId.Hex(),
	})
}

func (h *AdminController) UpdateModelOption(c *gin.Context) {
	ctx := context.Background()

	adminIdStr := c.Param("adminID")
	adminId, err := strconv.ParseUint(adminIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	modelOptionIdStr := c.Param("modelOptionID")
	objectId, err := primitive.ObjectIDFromHex(modelOptionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid model option ID"})
		return
	}

	var modelOption entity.ModelOption
	if err := c.ShouldBindJSON(&modelOption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	modelOption.ID = objectId
	modelOption.UpdatedAt = time.Now()

	if err := h.adminService.UpdateModelOption(ctx, adminId, modelOption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to update model option: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "model option updated successfully"})
}
