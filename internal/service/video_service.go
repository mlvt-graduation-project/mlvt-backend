package service

import (
	"errors"
	"log"
	"mlvt/internal/models"
	"mlvt/internal/repository"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type VideoService struct {
	videoRepo repository.VideoRepository
	awsRepo   repository.AWSRepository
}

func NewVideoService(videoRepo repository.VideoRepository, awsRepo repository.AWSRepository) *VideoService {
	return &VideoService{
		videoRepo: videoRepo,
		awsRepo:   awsRepo,
	}
}

func generateUUID() string {
	return uuid.New().String()
}

func getVideoDuration(output string) (float64, error) {
	durationIdx := strings.Index(output, "Duration: ")
	if durationIdx == -1 {
		return 0, errors.New("cannot find duration in FFmpeg output")
	}

	durationStr := output[durationIdx+10 : durationIdx+21] // expected format "00:00:00.00"
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

func (s *VideoService) ProcessVideo(filePath string) (*models.Video, error) {
	duration, fileSize, err := getVideoMetadata(filePath)
	if err != nil {
		return nil, err
	}

	video := &models.Video{
		ID:         generateUUID(),
		FilePath:   filePath,
		UploadedAt: time.Now(),
		Size:       fileSize,
		Duration:   int64(duration),
	}

	return video, nil
}

func (s *VideoService) UploadAndSaveVideo(video *models.Video, bucketName string) error {
	bucketKey := video.ID + ".mp4"
	err := s.awsRepo.UploadFile(bucketName, bucketKey, video.FilePath)
	if err != nil {
		return err
	}

	video.FilePath = bucketKey //"s3://" + bucketName + "/" + bucketKey
	video.UploadedAt = time.Now()

	log.Printf("Video metadata file uploaded - ID: %s, S3 Path: %s", video.ID, video.FilePath)
	return nil
}

func (s *VideoService) GetUserVideos(userID string) ([]models.Video, error) {
	return s.videoRepo.GetVideosByUserID(userID)
}
