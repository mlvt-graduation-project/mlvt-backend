package handler

import (
	"net/http"
	"strconv"

	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/service"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	videoService service.VideoService
}

func NewVideoController(videoService service.VideoService) *VideoController {
	return &VideoController{videoService: videoService}
}

// AddVideo handles adding a new video
// @Summary Add a new video
// @Description Creates a new video record in the system
// @Tags videos
// @Accept json
// @Produce json
// @Param video body entity.Video true "Video data"
// @Success 201 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /videos [post]
func (h *VideoController) AddVideo(c *gin.Context) {
	var video entity.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.videoService.CreateVideo(&video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Video added successfully"})
}

// GenerateUploadURLForVideo generates a presigned URL for uploading a video file
// @Summary Generate presigned upload URL for a video
// @Description Generates a presigned URL to upload a video file to S3
// @Tags videos
// @Produce json
// @Param file_name query string true "Name of the video file"
// @Param file_type query string true "Type of the video file (e.g., video/mp4)"
// @Success 200 {object} map[string]string "upload_url"
// @Failure 500 {object} map[string]string "error"
// @Router /videos/upload-url/video [post]
func (h *VideoController) GenerateUploadURLForVideo(c *gin.Context) {
	folder := env.EnvConfig.VideosFolder
	fileName := c.Query("file_name")
	fileType := c.Query("file_type")

	url, err := h.videoService.GeneratePresignedUploadURLForVideo(folder, fileName, fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url})
}

// GenerateUploadURLForImage generates a presigned URL for uploading an image file
// @Summary Generate presigned upload URL for an image
// @Description Generates a presigned URL to upload an image (e.g., thumbnail) to S3
// @Tags videos
// @Produce json
// @Param file_name query string true "Name of the image file"
// @Param file_type query string true "Type of the image file (e.g., image/jpeg)"
// @Success 200 {object} map[string]string "upload_url"
// @Failure 500 {object} map[string]string "error"
// @Router /videos/upload-url/image [post]
func (h *VideoController) GenerateUploadURLForImage(c *gin.Context) {
	folder := env.EnvConfig.VideoFramesFolder
	fileName := c.Query("file_name")
	fileType := c.Query("file_type")

	url, err := h.videoService.GeneratePresignedUploadURLForImage(folder, fileName, fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url})
}

// GenerateDownloadURLForVideo generates a presigned URL for downloading a video file
// @Summary Generate presigned download URL for a video
// @Description Generates a presigned URL to download a video file from S3
// @Tags videos
// @Produce json
// @Param video_id path int true "ID of the video file"
// @Success 200 {object} map[string]string "video_download_url"
// @Failure 500 {object} map[string]string "error"
// @Router /videos/{video_id}/download-url/video [get]
func (h *VideoController) GenerateDownloadURLForVideo(c *gin.Context) {
	videoID, _ := strconv.ParseUint(c.Param("video_id"), 10, 64)

	// Call the service to generate the presigned download URL for the video
	downloadURL, err := h.videoService.GeneratePresignedDownloadURLForVideo(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"video_download_url": downloadURL})
}

// GenerateDownloadURLForImage generates a presigned URL for downloading an image file
// @Summary Generate presigned download URL for an image
// @Description Generates a presigned URL to download an image (e.g., thumbnail) from S3
// @Tags videos
// @Produce json
// @Param video_id path int true "ID of the video file"
// @Success 200 {object} map[string]string "image_download_url"
// @Failure 500 {object} map[string]string "error"
// @Router /videos/{video_id}/download-url/image [get]
func (h *VideoController) GenerateDownloadURLForImage(c *gin.Context) {
	videoID, _ := strconv.ParseUint(c.Param("video_id"), 10, 64)

	// Call the service to generate the presigned download URL for the image
	downloadURL, err := h.videoService.GeneratePresignedDownloadURLForImage(videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image_download_url": downloadURL})
}

// GetVideoByID handles fetching a video by its ID and generating download URLs for the video and image
// @Summary Get video by ID
// @Description Fetches a video by its ID and generates presigned URLs for the video and image
// @Tags videos
// @Produce json
// @Param video_id path int true "ID of the video"
// @Success 200 {object} map[string]interface{} "video, video_url, image_url"
// @Failure 404 {object} map[string]string "error"
// @Router /videos/{video_id} [get]
func (h *VideoController) GetVideoByID(c *gin.Context) {
	videoID, _ := strconv.ParseUint(c.Param("video_id"), 10, 64)

	video, videoURL, imageURL, err := h.videoService.GetVideoByID(videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"video":     video,
		"video_url": videoURL,
		"image_url": imageURL,
	})
}

// DeleteVideo handles the deletion of a video by its ID
// @Summary Delete a video
// @Description Deletes a video by its ID from the system
// @Tags videos
// @Produce json
// @Param video_id path int true "ID of the video"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /videos/{video_id} [delete]
func (h *VideoController) DeleteVideo(c *gin.Context) {
	videoID, _ := strconv.ParseUint(c.Param("video_id"), 10, 64)

	if err := h.videoService.DeleteVideo(videoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
}

// ListVideosByUserID handles listing all videos for a specific user along with presigned image URLs
// @Summary List videos by user ID
// @Description Fetches all videos for a specific user along with presigned image URLs
// @Tags videos
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "videos, frames"
// @Failure 500 {object} map[string]string "error"
// @Router /videos/user/{user_id} [get]
func (h *VideoController) ListVideosByUserID(c *gin.Context) {
	// Parse user ID from the URL path
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	// Call the service to get the list of videos and presigned image URLs
	videos, frames, err := h.videoService.ListVideosByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the videos and the frames with presigned URLs
	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
		"frames": frames,
	})
}

// GetVideoStatus godoc
// @Summary Get the status of a video
// @Description Retrieve the status of a specific video by its ID
// @Tags Videos
// @Accept  json
// @Produce  json
// @Param   video_id path     uint64 true "Video ID"
// @Success 200 {object} map[string]string{"status": "success"}
// @Failure 400 {object} map[string]string{"error": "invalid video ID"}
// @Failure 404 {object} map[string]string{"error": "video not found"}
// @Failure 500 {object} map[string]string{"error": "internal server error"}
// @Router /videos/{video_id}/status [get]
func (h *VideoController) GetVideoStatus(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
		return
	}

	status, err := h.videoService.GetVideoStatus(videoID)

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// UpdateVideoStatusRequest represent the request body for updating video status
type UpdateVideoStatusRequest struct {
	Status entity.VideoStatus `json:"status" binding:"required, oneof=raw processing failed success`
}

// UpdateVideoStatus godoc
// @Summary Update the status of a video
// @Description Update the status of a specific video by its ID
// @Tags Videos
// @Accept  json
// @Produce  json
// @Param   video_id path     uint64 true "Video ID"
// @Param   status   body     UpdateVideoStatusRequest true "New status"
// @Success 200 {object} map[string]string{"message": "status updated successfully"}
// @Failure 400 {object} map[string]string{"error": "invalid input"}
// @Failure 404 {object} map[string]string{"error": "video not found"}
// @Failure 500 {object} map[string]string{"error": "internal server error"}
// @Router /videos/{video_id}/status [put]
func (vc *VideoController) UpdateVideoStatus(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
		return
	}

	var req UpdateVideoStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err = vc.videoService.UpdateVideoStatus(videoID, req.Status)
	if err != nil {
		if err.Error() == "no video found with id "+strconv.FormatUint(videoID, 10) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "video not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}
