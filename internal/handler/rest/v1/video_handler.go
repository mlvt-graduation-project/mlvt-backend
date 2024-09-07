package handler

import (
	"errors"
	"mlvt/internal/entity"
	"mlvt/internal/infra/reason"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/json"
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
//	@Summary		Add a new video
//	@Description	Add a new video for a specific user
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			video	body		schema.AddVideoRequest	true	"Video details"
//	@Success		201		{object}	map[string]string		"Video added successfully"
//	@Failure		400		{object}	map[string]string		"Invalid request"
//	@Failure		500		{object}	map[string]string		"Internal server error"
//	@Router			/videos [post]
func (vc *VideoController) AddVideo(ctx *gin.Context) {
	var req schema.AddVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		json.ErrorJSON(ctx, err, http.StatusBadRequest)
		return
	}

	err := vc.videoService.AddVideo(req.UserID, req.Title, req.Link, req.Duration)
	if err != nil {
		json.ErrorJSON(ctx, err, http.StatusInternalServerError)
		return
	}

	json.WriteJSON(ctx, http.StatusCreated, gin.H{"message": reason.VideoAdded.Message()})
}

// GetVideo fetches video details by ID
//	@Summary		Get video by ID
//	@Description	Get details of a video by its ID
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"Video ID"
//	@Success		200	{object}	schema.Video		"Video details"
//	@Failure		400	{object}	map[string]string	"Invalid video ID"
//	@Failure		404	{object}	map[string]string	"Video not found"
//	@Router			/videos/{id} [get]
func (vc *VideoController) GetVideo(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		json.ErrorJSON(ctx, errors.New(reason.InvalidVideoID.Message()), http.StatusBadRequest)
		return
	}

	video, err := vc.videoService.GetVideo(id)
	if err != nil {
		json.ErrorJSON(ctx, errors.New(reason.VideoNotFound.Message()), http.StatusNotFound)
		return
	}

	// Convert entity.Video to schema.Video for response
	schemaVideo := convertEntityVideoToSchemaVideo(*video)

	json.WriteJSON(ctx, http.StatusOK, schemaVideo)
}

// UpdateVideo handles updating video details
//	@Summary		Update video details
//	@Description	Update the details of an existing video
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Video ID"
//	@Param			video	body		schema.AddVideoRequest	true	"Updated video details"
//	@Success		200		{object}	map[string]string		"Video updated successfully"
//	@Failure		400		{object}	map[string]string		"Invalid video ID or request"
//	@Failure		500		{object}	map[string]string		"Internal server error"
//	@Router			/videos/{id} [put]
func (vc *VideoController) UpdateVideo(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		json.ErrorJSON(ctx, errors.New(reason.InvalidVideoID.Message()), http.StatusBadRequest)
		return
	}

	var req schema.AddVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		json.ErrorJSON(ctx, err, http.StatusBadRequest)
		return
	}

	err = vc.videoService.UpdateVideo(id, req.Link, req.Duration)
	if err != nil {
		json.ErrorJSON(ctx, err, http.StatusInternalServerError)
		return
	}

	json.WriteJSON(ctx, http.StatusOK, map[string]string{"message": reason.VideoUpdated.Message()})

}

// DeleteVideo handles deleting a video by ID
//	@Summary		Delete video
//	@Description	Delete a video by its ID
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"Video ID"
//	@Success		200	{object}	map[string]string	"Video deleted successfully"
//	@Failure		400	{object}	map[string]string	"Invalid video ID"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/videos/{id} [delete]
func (vc *VideoController) DeleteVideo(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		json.ErrorJSON(ctx, errors.New(reason.InvalidVideoID.Message()), http.StatusBadRequest)
		return
	}

	err = vc.videoService.DeleteVideo(id)
	if err != nil {
		json.ErrorJSON(ctx, err, http.StatusInternalServerError)
		return
	}

	json.WriteJSON(ctx, http.StatusOK, map[string]string{"message": reason.VideoDeleted.Message()})
}

// GetVideosByUser handles fetching all videos for a specific user
//	@Summary		Get videos by user ID
//	@Description	Get all videos uploaded by a specific user
//	@Tags			videos
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int							true	"User ID"
//	@Success		200		{object}	schema.GetVideosResponse	"List of videos"
//	@Failure		400		{object}	map[string]string			"Invalid user ID"
//	@Failure		500		{object}	map[string]string			"Internal server error"
//	@Router			/videos/user/{userID} [get]
func (vc *VideoController) GetVideosByUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 64)
	if err != nil {
		json.ErrorJSON(ctx, errors.New(reason.InvalidUserID.Message()), http.StatusBadRequest)
		return
	}

	videos, err := vc.videoService.GetVideosByUser(userID)
	if err != nil {
		json.ErrorJSON(ctx, err, http.StatusInternalServerError)
		return
	}

	// Convert []entity.Video to []schema.Video
	var schemaVideos []schema.Video
	for _, v := range videos {
		schemaVideos = append(schemaVideos, convertEntityVideoToSchemaVideo(v))
	}

	json.WriteJSON(ctx, http.StatusOK, schema.GetVideosResponse{Videos: schemaVideos})
}

// GeneratePresignedURL handles the request to generate a pre-signed URL for video uploads
//	@Summary		Generate a pre-signed URL for video uploads
//	@Description	Generates a pre-signed URL for uploading a video to AWS S3 and registers the video with initial data
//	@Tags			Videos
//	@Accept			json
//	@Produce		json
//	@Param			request	body		schema.GeneratePresignedURLRequest	true	"Video upload details"
//	@Success		200		{object}	schema.GeneratePresignedURLResponse
//	@Failure		400		{object}	schema.ErrorResponse	"Invalid request format"
//	@Failure		500		{object}	schema.ErrorResponse	"Failed to generate presigned URL or Failed to add video"
//	@Router			/videos/generate-presigned-url [post]
func (vc *VideoController) GeneratePresignedURLHandler(ctx *gin.Context) {
	// Bind the JSON payload to the request struct and validate it
	var req schema.GeneratePresignedURLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		json.WriteJSON(ctx, http.StatusBadRequest, schema.ErrorResponse{Error: reason.InvalidRequestFormat.Message()})
		return
	}

	// Generate a pre-signed URL for the video upload using the service layer
	url, err := vc.videoService.GeneratePresignedURL(req.FileName, req.FileType)
	if err != nil {
		json.WriteJSON(ctx, http.StatusInternalServerError, schema.ErrorResponse{Error: reason.FailedToGeneratePresignedURL.Message()})
		return
	}

	// Use the AddVideo method to register the video with initial data
	err = vc.videoService.AddVideo(req.UserID, req.Title, url, req.Duration)
	if err != nil {
		log.Errorf(reason.FailedToAddVideo.Message(), err)
		json.WriteJSON(ctx, http.StatusInternalServerError, schema.ErrorResponse{Error: reason.FailedToAddVideo.Message()})
		return
	}

	// Send the presigned URL back to the frontend
	json.WriteJSON(ctx, http.StatusOK, schema.GeneratePresignedURLResponse{PresignedUrl: url})
}

// Helper function to convert entity.Video to schema.Video
func convertEntityVideoToSchemaVideo(entityVideo entity.Video) schema.Video {
	return schema.Video{
		ID:        entityVideo.ID,
		Title:     entityVideo.Title,
		Duration:  entityVideo.Duration,
		Link:      entityVideo.Link,
		UserID:    entityVideo.UserID,
		CreatedAt: entityVideo.CreatedAt,
		UpdatedAt: entityVideo.UpdatedAt,
	}
}
