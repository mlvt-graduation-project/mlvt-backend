package media_service

import (
	"fmt"
	"mlvt/internal/entity"
)

func (s *mediaService) CreateTranscription(transcription *entity.Transcription) (uint64, error) {
	return s.mediaRepo.CreateTranscription(transcription)
}

func (s *mediaService) GetTranscriptionByID(transcriptionID uint64) (*entity.Transcription, string, error) {
	transcription, err := s.mediaRepo.GetTranscriptionByID(transcriptionID)
	if err != nil {
		return nil, "", err
	}
	if transcription == nil {
		return nil, "", fmt.Errorf("transcription not found")
	}

	// Generate presigned URL
	presignedURL, err := s.s3Client.GeneratePresignedDownloadURL(transcription.Folder, transcription.FileName, "text/plain")
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate presigned download URL: %v", err)
	}

	return transcription, presignedURL, nil
}

func (s *mediaService) GetTranscriptionByIDAndUserID(transcriptionID, userID uint64) (*entity.Transcription, string, error) {
	transcription, err := s.mediaRepo.GetTranscriptionByIDAndUserID(transcriptionID, userID)
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

func (s *mediaService) GetTranscriptionByIDAndVideoID(transcriptionID, videoID uint64) (*entity.Transcription, string, error) {
	transcription, err := s.mediaRepo.GetTranscriptionByIDAndVideoID(transcriptionID, videoID)
	if err != nil {
		return nil, "", err
	}
	if transcription == nil {
		return nil, "", fmt.Errorf("transcription not found")
	}

	// Generate presigned URL
	presignedURL, err := s.s3Client.GeneratePresignedDownloadURL(transcription.Folder, transcription.FileName, "text/plain")
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate presigned download URL: %v", err)
	}

	return transcription, presignedURL, nil
}

func (s *mediaService) ListTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error) {
	return s.mediaRepo.ListTranscriptionsByUserID(userID)
}

func (s *mediaService) ListTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error) {
	return s.mediaRepo.ListTranscriptionsByVideoID(videoID)
}

func (s *mediaService) DeleteTranscription(transcriptionID uint64) error {
	return s.mediaRepo.DeleteTranscription(transcriptionID)
}

func (s *mediaService) GeneratePresignedUploadURLForText(folder, fileName, fileType string) (string, error) {
	return s.s3Client.GeneratePresignedURL(folder, fileName, fileType)
}

func (s *mediaService) GeneratePresignedDownloadURLForText(transcriptionID uint64) (string, error) {
	transcription, err := s.mediaRepo.GetTranscriptionByID(transcriptionID)
	if err != nil {
		return "", err
	}
	if transcription == nil {
		return "", fmt.Errorf("transcription not found")
	}

	return s.s3Client.GeneratePresignedDownloadURL(transcription.Folder, transcription.FileName, "text/plain")
}

func (s *mediaService) UpdateTranscription(transcription *entity.Transcription) error {
	return s.mediaRepo.UpdateTranscription(transcription)
}

func (s *mediaService) UpdateTranscriptionStatus(transcriptionID uint64, status entity.StatusEntity) error {
	return s.mediaRepo.UpdateTranscriptionStatus(transcriptionID, status)
}
