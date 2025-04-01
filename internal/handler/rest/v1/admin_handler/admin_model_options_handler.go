package admin_handler

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	if modelOption.ID == primitive.NilObjectID {
		modelOption.ID = primitive.NewObjectID()
	}
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
