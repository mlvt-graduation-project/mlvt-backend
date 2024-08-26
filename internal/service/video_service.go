package service

import (
	"database/sql"
	"errors"
	"mlvt/internal/entity"
	"mlvt/internal/repo"
)

// VideoService handles video-related business logic
type VideoService struct {
	videoRepo repo.VideoRepository
}

// NewVideoService creates a new instance of VideoService
func NewVideoService(videoRepo repo.VideoRepository) *VideoService {
	return &VideoService{
		videoRepo: videoRepo,
	}
}

// AddVideo adds a new video for a user
func (s *VideoService) AddVideo(userID uint64, link string, duration int) error {
	if link == "" {
		return errors.New("video link cannot be empty")
	}
	if duration <= 0 {
		return errors.New("video duration must be positive")
	}

	video := &entity.Video{
		UserID:   userID,
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
			return nil, errors.New("video not found")
		}
		return nil, err
	}
	return video, nil
}

// UpdateVideo updates an existing video's details
func (s *VideoService) UpdateVideo(id uint64, link string, duration int) error {
	if link == "" {
		return errors.New("video link cannot be empty")
	}
	if duration <= 0 {
		return errors.New("video duration must be positive")
	}

	video, err := s.videoRepo.GetVideoByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("video not found")
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
			return nil, errors.New("no videos found for user")
		}
		return nil, err
	}
	return videos, nil
}
