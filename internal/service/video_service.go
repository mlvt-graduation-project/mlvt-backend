package service

import (
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/infra/env"
	"mlvt/internal/pkg/response"
	"mlvt/internal/repo"
)

type VideoService interface {
	CreateVideo(video *entity.Video) (uint64, error)
	GetVideoByID(videoID uint64) (*entity.Video, string, string, error) // Returns the video record and presigned URLs for video and image
	ListVideosByUserID(userID uint64) ([]response.ListVideosByUserIDResponse, error)
	DeleteVideo(videoID uint64) error
	UpdateVideo(video *entity.Video) error
	UpdateVideoStatus(videoID uint64, status entity.StatusEntity) error
	GetVideoStatus(videoID uint64) (entity.StatusEntity, error)
	GeneratePresignedUploadURLForVideo(folder, fileName, fileType string) (string, error)
	GeneratePresignedUploadURLForImage(folder, fileName, fileType string) (string, error)
	GeneratePresignedDownloadURLForVideo(videoID uint64) (string, error)
	GeneratePresignedDownloadURLForImage(videoID uint64) (string, error)
}

type videoService struct {
	repo     repo.VideoRepository
	s3Client aws.S3ClientInterface
}

func NewVideoService(repo repo.VideoRepository, s3Client aws.S3ClientInterface) VideoService {
	return &videoService{
		repo:     repo,
		s3Client: s3Client,
	}
}

func (s *videoService) CreateVideo(video *entity.Video) (uint64, error) {
	id, err := s.repo.CreateVideo(video)
	return id, err
}

func (s *videoService) GetVideoByID(videoID uint64) (*entity.Video, string, string, error) {
	video, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return nil, "", "", err
	}
	if video == nil {
		return nil, "", "", fmt.Errorf("video not found")
	}

	// Generate presigned URLs for video and image
	videoURL, err := s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.FileName, "video/mp4")
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate presigned video URL: %v", err)
	}
	imageURL, err := s.s3Client.GeneratePresignedDownloadURL(env.EnvConfig.VideoFramesFolder, video.Image, "image/jpeg")
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate presigned image URL: %v", err)
	}

	return video, videoURL, imageURL, nil
}

func (s *videoService) ListVideosByUserID(userID uint64) ([]response.ListVideosByUserIDResponse, error) {
	// Fetch the videos for the user
	videos, err := s.repo.ListVideosByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list videos for user %d: %v", userID, err)
	}

	var videoWithURLsList []response.ListVideosByUserIDResponse
	for _, video := range videos {
		// Generate the presigned URL for the video
		videoURL, err := s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.FileName, "video/mp4")
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned video URL for video ID %d: %v", video.ID, err)
		}

		// Generate the presigned URL for the image/frame
		imageURL, err := s.s3Client.GeneratePresignedDownloadURL(env.EnvConfig.VideoFramesFolder, video.Image, "image/jpeg")
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned image URL for video ID %d: %v", video.ID, err)
		}

		// Append to the list
		videoWithURLsList = append(videoWithURLsList, response.ListVideosByUserIDResponse{
			Video:    video,
			VideoURL: videoURL,
			ImageURL: imageURL,
		})
	}

	return videoWithURLsList, nil
}

func (s *videoService) DeleteVideo(videoID uint64) error {
	// Fetch the video record to get the file names
	video, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return fmt.Errorf("failed to fetch video: %v", err)
	}
	if video == nil {
		return fmt.Errorf("video not found")
	}

	// Begin deletion process
	// 1. Delete the video and frame files from S3
	err = s.s3Client.DeleteFile(video.Folder, video.FileName)
	if err != nil {
		return fmt.Errorf("failed to delete video file from S3: %v", err)
	}

	err = s.s3Client.DeleteFile(env.EnvConfig.VideoFramesFolder, video.Image)
	if err != nil {
		return fmt.Errorf("failed to delete frame image from S3: %v", err)
	}

	// 2. Delete the video record from the database
	err = s.repo.DeleteVideo(videoID)
	if err != nil {
		return fmt.Errorf("failed to delete video from database: %v", err)
	}

	return nil
}

func (s *videoService) UpdateVideo(video *entity.Video) error {
	return s.repo.UpdateVideo(video)
}

func (s *videoService) UpdateVideoStatus(videoID uint64, status entity.StatusEntity) error {
	return s.repo.UpdateVideoStatus(videoID, status)
}
func (s *videoService) GetVideoStatus(videoID uint64) (entity.StatusEntity, error) {
	return s.repo.GetVideoStatus(videoID)
}

// GeneratePresignedUploadURLForVideo generates a presigned URL for uploading a video file
func (s *videoService) GeneratePresignedUploadURLForVideo(folder, fileName, fileType string) (string, error) {
	return s.s3Client.GeneratePresignedURL(folder, fileName, fileType)
}

// GeneratePresignedUploadURLForImage generates a presigned URL for uploading an image file
func (s *videoService) GeneratePresignedUploadURLForImage(folder, fileName, fileType string) (string, error) {
	return s.s3Client.GeneratePresignedURL(folder, fileName, fileType)
}

// GeneratePresignedDownloadURLForVideo generates a presigned URL for downloading a video file
func (s *videoService) GeneratePresignedDownloadURLForVideo(videoID uint64) (string, error) {
	video, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return "", err
	}
	if video == nil {
		return "", fmt.Errorf("video not found")
	}

	return s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.FileName, "video/mp4")
}

// GeneratePresignedDownloadURLForImage generates a presigned URL for downloading an image file
func (s *videoService) GeneratePresignedDownloadURLForImage(videoID uint64) (string, error) {
	video, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return "", err
	}
	if video == nil {
		return "", fmt.Errorf("video not found")
	}

	return s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.Image, "image/jpeg")
}
