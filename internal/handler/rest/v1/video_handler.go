package handler

import (
	"mlvt/internal/entity"
	"mlvt/internal/schema"
	"mlvt/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// VideoController manages video-related requests
type VideoController struct {
	videoService *service.VideoService
}

// NewVideoController creates a new VideoController
func NewVideoController(videoService *service.VideoService) *VideoController {
	return &VideoController{
		videoService: videoService,
	}
}

// AddVideo handles adding a new video
// @Summary Add a new video
// @Description Add a new video for a specific user
// @Tags videos
// @Accept json
// @Produce json
// @Param video body schema.AddVideoRequest true "Video details"
// @Success 201 {object} map[string]string "Video added successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /videos [post]
func (vc *VideoController) AddVideo(ctx *gin.Context) {
	var req schema.AddVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := vc.videoService.AddVideo(req.UserID, req.Link, req.Duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Video added successfully"})
}

// GetVideo fetches video details by ID
// @Summary Get video by ID
// @Description Get details of a video by its ID
// @Tags videos
// @Accept json
// @Produce json
// @Param id path int true "Video ID"
// @Success 200 {object} schema.Video "Video details"
// @Failure 400 {object} map[string]string "Invalid video ID"
// @Failure 404 {object} map[string]string "Video not found"
// @Router /videos/{id} [get]
func (vc *VideoController) GetVideo(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	video, err := vc.videoService.GetVideo(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Convert entity.Video to schema.Video for response
	schemaVideo := convertEntityVideoToSchemaVideo(*video)

	ctx.JSON(http.StatusOK, schemaVideo)
}

// UpdateVideo handles updating video details
// @Summary Update video details
// @Description Update the details of an existing video
// @Tags videos
// @Accept json
// @Produce json
// @Param id path int true "Video ID"
// @Param video body schema.AddVideoRequest true "Updated video details"
// @Success 200 {object} map[string]string "Video updated successfully"
// @Failure 400 {object} map[string]string "Invalid video ID or request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /videos/{id} [put]
func (vc *VideoController) UpdateVideo(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	var req schema.AddVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = vc.videoService.UpdateVideo(id, req.Link, req.Duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Video updated successfully"})
}

// DeleteVideo handles deleting a video by ID
// @Summary Delete video
// @Description Delete a video by its ID
// @Tags videos
// @Accept json
// @Produce json
// @Param id path int true "Video ID"
// @Success 200 {object} map[string]string "Video deleted successfully"
// @Failure 400 {object} map[string]string "Invalid video ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /videos/{id} [delete]
func (vc *VideoController) DeleteVideo(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	err = vc.videoService.DeleteVideo(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
}

// GetVideosByUser handles fetching all videos for a specific user
// @Summary Get videos by user ID
// @Description Get all videos uploaded by a specific user
// @Tags videos
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {object} schema.GetVideosResponse "List of videos"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /videos/user/{userID} [get]
func (vc *VideoController) GetVideosByUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	videos, err := vc.videoService.GetVideosByUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert []entity.Video to []schema.Video
	var schemaVideos []schema.Video
	for _, v := range videos {
		schemaVideos = append(schemaVideos, convertEntityVideoToSchemaVideo(v))
	}

	ctx.JSON(http.StatusOK, schema.GetVideosResponse{Videos: schemaVideos})
}

// Helper function to convert entity.Video to schema.Video
func convertEntityVideoToSchemaVideo(entityVideo entity.Video) schema.Video {
	return schema.Video{
		ID:        entityVideo.ID,
		Duration:  entityVideo.Duration,
		Link:      entityVideo.Link,
		UserID:    entityVideo.UserID,
		CreatedAt: entityVideo.CreatedAt,
		UpdatedAt: entityVideo.UpdatedAt,
	}
}
