package handlers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"mlvt/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func generateUUID() string {
	return uuid.New().String()
}

func getVideoDuration(output string) (float64, error) {
	durationIdx := strings.Index(output, "Duration: ")
	if durationIdx == -1 {
		return 0, errors.New("cannot find duration in FFmpeg output")
	}

	// Extracting the duration part of the output
	durationStr := output[durationIdx+10 : durationIdx+21] // Expected format "00:00:00.00"
	timeParts := strings.Split(durationStr, ":")
	if len(timeParts) != 3 {
		return 0, errors.New("unexpected duration format")
	}

	hours, err := strconv.Atoi(strings.TrimSpace(timeParts[0]))
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(strings.TrimSpace(timeParts[1]))
	if err != nil {
		return 0, err
	}
	secondsParts := strings.Split(timeParts[2], ".")
	seconds, err := strconv.Atoi(strings.TrimSpace(secondsParts[0]))
	if err != nil {
		return 0, err
	}

	return float64(hours*3600 + minutes*60 + seconds), nil
}

func getVideoMetadata(filePath string) (duration float64, fileSize int64, err error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, 0, err
	}
	fileSize = fileInfo.Size()

	cmd := exec.Command("ffmpeg", "-i", filePath, "-f", "null", "-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("FFmpeg output: \n%s", output)
		return 0, fileSize, err
	}

	duration, err = getVideoDuration(string(output))
	if err != nil {
		return 0, fileSize, err
	}

	return duration, fileSize, nil
}

func UploadVideos(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
		return
	}

	// Temporarily save the file to process
	tempPath := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video for processing"})
		return
	}
	defer os.Remove(tempPath) // Clean up after processing

	duration, fileSize, err := getVideoMetadata(tempPath)
	if err != nil {
		log.Printf("Error processing video: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process video"})
		return
	}

	video := models.Video{
		ID:         generateUUID(),
		FilePath:   tempPath,
		UploadedAt: time.Now(),
		Size:       fileSize,
		Duration:   int64(duration),
		Type:       file.Header.Get("Content-Type"),
	}

	log.Printf("Received video - ID: %s, Duration: %d sec, Type: %s, UploadedAt: %s, Size: %d bytes, UserID: %s, Type: %s",
		video.ID, video.Duration, video.FilePath, video.UploadedAt.Format(time.RFC3339), video.Size, video.UserID, video.Type)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Video uploaded successfully",
		"data": gin.H{
			"id":         video.ID,
			"duration":   video.Duration,
			"uploadedAt": video.UploadedAt.Format(time.RFC3339),
			"size":       video.Size,
			"type":       video.Type,
		},
	})
}
