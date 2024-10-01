package service

import (
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/repo"
)

type AudioService interface {
	GeneratePresignedUploadURL(folder, fileName, fileType string) (string, error)
	GeneratePresignedDownloadURL(audioID uint64) (string, error)
	CreateAudio(audio *entity.Audio) error
	GetAudioByID(audioID uint64) (*entity.Audio, string, error)
	GetAudioByIDAndUserID(audioID, userID uint64) (*entity.Audio, string, error)
	ListAudiosByUserID(userID uint64) ([]entity.Audio, error)
	GetAudioByVideoID(videoID, audioID uint64) (*entity.Audio, string, error)
	ListAudiosByVideoID(videoID uint64) ([]entity.Audio, error)
	DeleteAudio(audioID uint64) error
}

type audioService struct {
	repo     repo.AudioRepository
	s3Client *aws.S3Client
}

func NewAudioService(repo repo.AudioRepository, s3Client *aws.S3Client) AudioService {
	return &audioService{
		repo:     repo,
		s3Client: s3Client,
	}
}

func (s *audioService) GeneratePresignedUploadURL(folder, fileName, fileType string) (string, error) {
	return s.s3Client.GeneratePresignedURL(folder, fileName, fileType)
}

func (s *audioService) GeneratePresignedDownloadURL(audioID uint64) (string, error) {
	// Fetch the audio from the repository using its ID
	audio, err := s.repo.GetAudioByID(audioID)
	if err != nil {
		return "", fmt.Errorf("could not find audio with ID %d: %v", audioID, err)
	}

	// Generate the presigned URL using S3 client
	presignedURL, err := s.s3Client.GeneratePresignedURL(audio.Folder, audio.FileName, "audio/mpeg")
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %v", err)
	}

	// Return the generated presigned URL
	return presignedURL, nil
}

func (s *audioService) CreateAudio(audio *entity.Audio) error {
	return s.repo.CreateAudio(audio)
}

func (s *audioService) GetAudioByID(audioID uint64) (*entity.Audio, string, error) {
	audio, err := s.repo.GetAudioByID(audioID)
	if err != nil {
		return nil, "", err
	}
	presignedURL, err := s.s3Client.GeneratePresignedURL(audio.Folder, audio.FileName, "audio/mpeg")
	if err != nil {
		return nil, "", err
	}
	return audio, presignedURL, nil
}

// GetAudioByIDAndUserID retrieves a single audio by its ID and User ID and generates a presigned URL
func (s *audioService) GetAudioByIDAndUserID(audioID, userID uint64) (*entity.Audio, string, error) {
	// Fetch the audio from the repository
	audio, err := s.repo.GetAudioByIDAndUserID(audioID, userID)
	if err != nil {
		return nil, "", err
	}
	if audio == nil {
		return nil, "", fmt.Errorf("audio not found")
	}

	// Generate the presigned URL using the S3 client
	presignedURL, err := s.s3Client.GeneratePresignedURL(audio.Folder, audio.FileName, "audio/mpeg")
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate presigned download URL: %v", err)
	}

	return audio, presignedURL, nil
}
func (s *audioService) ListAudiosByUserID(userID uint64) ([]entity.Audio, error) {
	return s.repo.ListAudiosByUserID(userID)
}

func (s *audioService) GetAudioByVideoID(videoID, audioID uint64) (*entity.Audio, string, error) {
	audio, err := s.repo.GetAudioByVideoID(videoID, audioID)
	if err != nil {
		return nil, "", err
	}
	presignedURL, err := s.s3Client.GeneratePresignedURL(audio.Folder, audio.FileName, "audio/mpeg")
	if err != nil {
		return nil, "", err
	}
	return audio, presignedURL, nil
}

func (s *audioService) ListAudiosByVideoID(videoID uint64) ([]entity.Audio, error) {
	return s.repo.ListAudiosByVideoID(videoID)
}

func (s *audioService) DeleteAudio(audioID uint64) error {
	return s.repo.DeleteAudioByID(audioID)
}
