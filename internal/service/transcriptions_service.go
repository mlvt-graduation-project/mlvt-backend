package service

import (
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/repo"
)

type TranscriptionService interface {
	GeneratePresignedURLForTranscription(videoID uint64) (string, error)
	CreateTranscription(videoID, userID uint64, lang, transcriptionText string) (entity.Transcription, error)
	GetTranscription(userID, transcriptionID uint64) (entity.Transcription, error)
	ListTranscriptionsByUser(userID uint64) ([]entity.Transcription, error)
	GetTranscriptionByVideo(videoID, transcriptionID uint64) (entity.Transcription, error)
	ListTranscriptionsByVideo(videoID uint64) ([]entity.Transcription, error)
	DeleteTranscription(transcriptionID uint64) error
}

type transcriptionService struct {
	repo      repo.TranscriptionRepository
	s3Client  aws.S3Client
	videoRepo repo.VideoRepository
}

func NewTranscriptionService(repo repo.TranscriptionRepository, s3 aws.S3Client, videoRepo repo.VideoRepository) *transcriptionService {
	return &transcriptionService{
		repo:      repo,
		s3Client:  s3,
		videoRepo: videoRepo,
	}
}

// GeneratePresignedURLForTranscription generates a presigned URL for a given video
func (s *transcriptionService) GeneratePresignedURLForTranscription(videoID uint64) (string, error) {
	video, err := s.videoRepo.GetVideoByID(videoID)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve video: %w", err)
	}
	// Generate presigned URL
	fileType := "text/plain" // Default file type for transcription
	return s.s3Client.GeneratePresignedURL(video.Folder, video.FileName, fileType)
}

// CreateTranscription handles the logic to create and store a new transcription
func (s *transcriptionService) CreateTranscription(videoID, userID uint64, lang, transcriptionText string) (entity.Transcription, error) {
	video, err := s.videoRepo.GetVideoByID(videoID)
	if err != nil {
		return entity.Transcription{}, fmt.Errorf("unable to retrieve video: %w", err)
	}

	tx := entity.Transcription{
		VideoID:  videoID,
		UserID:   userID,
		Text:     transcriptionText,
		Lang:     lang,
		Folder:   video.Folder,
		FileName: video.FileName,
	}
	return s.repo.CreateTranscription(tx)
}

func (s *transcriptionService) GetTranscription(userID, transcriptionID uint64) (entity.Transcription, error) {
	return s.repo.GetTranscriptionByID(userID, transcriptionID)
}

func (s *transcriptionService) ListTranscriptionsByUser(userID uint64) ([]entity.Transcription, error) {
	return s.repo.GetTranscriptionsByUserID(userID)
}

func (s *transcriptionService) GetTranscriptionByVideo(videoID, transcriptionID uint64) (entity.Transcription, error) {
	return s.repo.GetTranscriptionByVideoID(videoID, transcriptionID)
}

func (s *transcriptionService) ListTranscriptionsByVideo(videoID uint64) ([]entity.Transcription, error) {
	return s.repo.GetTranscriptionsByVideoID(videoID)
}

func (s *transcriptionService) DeleteTranscription(transcriptionID uint64) error {
	return s.repo.DeleteTranscription(transcriptionID)
}
