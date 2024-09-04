package service

import (
	"database/sql"
	"errors"
	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/infra/reason"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/repo"
)

// VideoService handles video-related business logic
type VideoService struct {
	s3Client  *aws.S3Client
	videoRepo repo.VideoRepository
}

// NewVideoService creates a new instance of VideoService
func NewVideoService(s3Client *aws.S3Client, videoRepo repo.VideoRepository) *VideoService {
	return &VideoService{

		s3Client:  s3Client,
		videoRepo: videoRepo,
	}
}

// AddVideo adds a new video for a user
func (s *VideoService) AddVideo(userID uint64, title string, link string, duration int) error {
	if link == "" {
		return errors.New(reason.VideoLinkCannotBeEmpty.Message())
	}
	if duration <= 0 {
		return errors.New(reason.VideoDurationMustBePositive.Message())
	}
	if title == "" {
		return errors.New(reason.VideoTitleCannotBeEmpty.Message())
	}

	video := &entity.Video{
		UserID:   userID,
		Title:    title,
		Link:     link,
		Duration: duration,
	}
	return s.videoRepo.CreateVideo(video)
}

// GetVideo retrieves a video by ID
func (s *VideoService) GetVideo(id uint64) (*entity.Video, error) {
	video, err := s.videoRepo.GetVideoByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(reason.VideoNotFound.Message())
		}
		return nil, err
	}
	return video, nil
}

// UpdateVideo updates an existing video's details
func (s *VideoService) UpdateVideo(id uint64, link string, duration int) error {
	if link == "" {
		return errors.New(reason.VideoLinkCannotBeEmpty.Message())
	}
	if duration <= 0 {
		return errors.New(reason.VideoDurationMustBePositive.Message())
	}

	video, err := s.videoRepo.GetVideoByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(reason.VideoNotFound.Message())
		}
		return err
	}

	video.Link = link
	video.Duration = duration

	return s.videoRepo.UpdateVideo(video)
}

// DeleteVideo removes a video
func (s *VideoService) DeleteVideo(id uint64) error {
	return s.videoRepo.DeleteVideo(id)
}

// GetVideosByUser fetches all videos for a specific user
func (s *VideoService) GetVideosByUser(userID uint64) ([]entity.Video, error) {
	videos, err := s.videoRepo.GetVideosByUserID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(reason.NoVideoForUser.Message())
		}
		return nil, err
	}
	return videos, nil
}

// GeneratePresignedURL generates a pre-signed URL for uploading a video
func (s *VideoService) GeneratePresignedURL(fileName string, fileType string) (string, error) {
	log.Infof(reason.GeneratedPresignedURLForFile.Message() + ": " + fileName + "," + reason.Type.Message() + ":" + fileType)
	url, err := s.s3Client.GeneratePresignedURL(fileName, fileType)
	if err != nil {
		log.Errorf(reason.FailedToCreatePresignedURL.Message()+": %v", err)
		return "", err
	}
	log.Infof(reason.GeneratedPresignedURL.Message()+": %s", url)
	return url, nil
}
