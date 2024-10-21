package service

import (
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/repo"
)

type TranscriptionService interface {
	CreateTranscription(transcription *entity.Transcription) error
	GetTranscriptionByID(transcriptionID uint64) (*entity.Transcription, string, error)
	GetTranscriptionByIDAndUserID(transcriptionID, userID uint64) (*entity.Transcription, string, error)
	GetTranscriptionByIDAndVideoID(transcriptionID, videoID uint64) (*entity.Transcription, string, error)
	ListTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error)
	ListTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error)
	DeleteTranscription(transcriptionID uint64) error
	GeneratePresignedUploadURL(folder, fileName, fileType string) (string, error)
	GeneratePresignedDownloadURL(transcriptionID uint64) (string, error)
}

type transcriptionService struct {
	repo     repo.TranscriptionRepository
	s3Client aws.S3ClientInterface
}

func NewTranscriptionService(repo repo.TranscriptionRepository, s3Client aws.S3ClientInterface) TranscriptionService {
	return &transcriptionService{
		repo:     repo,
		s3Client: s3Client,
	}
}

func (s *transcriptionService) CreateTranscription(transcription *entity.Transcription) error {
	return s.repo.CreateTranscription(transcription)
}

func (s *transcriptionService) GetTranscriptionByID(transcriptionID uint64) (*entity.Transcription, string, error) {
	transcription, err := s.repo.GetTranscriptionByID(transcriptionID)
	if err != nil {
		return nil, "", err
	}
	if transcription == nil {
		return nil, "", fmt.Errorf("transcription not found")
	}

	// Generate presigned URL
	presignedURL, err := s.s3Client.GeneratePresignedURL(transcription.Folder, transcription.FileName, "application/json")
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate presigned download URL: %v", err)
	}

	return transcription, presignedURL, nil
}

func (s *transcriptionService) GetTranscriptionByIDAndUserID(transcriptionID, userID uint64) (*entity.Transcription, string, error) {
	transcription, err := s.repo.GetTranscriptionByIDAndUserID(transcriptionID, userID)
	if err != nil {
		return nil, "", err
	}
	if transcription == nil {
		return nil, "", fmt.Errorf("transcription not found")
	}

	// Generate presigned URL
	presignedURL, err := s.s3Client.GeneratePresignedURL(transcription.Folder, transcription.FileName, "application/json")
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate presigned download URL: %v", err)
	}

	return transcription, presignedURL, nil
}

func (s *transcriptionService) GetTranscriptionByIDAndVideoID(transcriptionID, videoID uint64) (*entity.Transcription, string, error) {
	transcription, err := s.repo.GetTranscriptionByIDAndVideoID(transcriptionID, videoID)
	if err != nil {
		return nil, "", err
	}
	if transcription == nil {
		return nil, "", fmt.Errorf("transcription not found")
	}

	// Generate presigned URL
	presignedURL, err := s.s3Client.GeneratePresignedURL(transcription.Folder, transcription.FileName, "application/json")
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate presigned download URL: %v", err)
	}

	return transcription, presignedURL, nil
}

func (s *transcriptionService) ListTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error) {
	return s.repo.ListTranscriptionsByUserID(userID)
}

func (s *transcriptionService) ListTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error) {
	return s.repo.ListTranscriptionsByVideoID(videoID)
}

func (s *transcriptionService) DeleteTranscription(transcriptionID uint64) error {
	return s.repo.DeleteTranscription(transcriptionID)
}

func (s *transcriptionService) GeneratePresignedUploadURL(folder, fileName, fileType string) (string, error) {
	return s.s3Client.GeneratePresignedURL(folder, fileName, fileType)
}

func (s *transcriptionService) GeneratePresignedDownloadURL(transcriptionID uint64) (string, error) {
	transcription, err := s.repo.GetTranscriptionByID(transcriptionID)
	if err != nil {
		return "", err
	}
	if transcription == nil {
		return "", fmt.Errorf("transcription not found")
	}

	return s.s3Client.GeneratePresignedURL(transcription.Folder, transcription.FileName, "application/json")
}
