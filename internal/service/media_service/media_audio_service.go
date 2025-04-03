package media_service

import (
	"fmt"
	"mlvt/internal/entity"
)

func (s *mediaService) GeneratePresignedUploadURL(folder, fileName, fileType string) (string, error) {
	return s.s3Client.GeneratePresignedURL(folder, fileName, fileType)
}

func (s *mediaService) GeneratePresignedDownloadURL(audioID uint64) (string, error) {
	// Fetch the audio from the repository using its ID
	audio, err := s.audioRepo.GetAudioByID(audioID)
	if err != nil {
		return "", fmt.Errorf("could not find audio with ID %d: %v", audioID, err)
	}

	// Generate the presigned URL using S3 client
	presignedURL, err := s.s3Client.GeneratePresignedDownloadURL(audio.Folder, audio.FileName, "audio/mpeg")
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL: %v", err)
	}

	// Return the generated presigned URL
	return presignedURL, nil
}

func (s *mediaService) CreateAudio(audio *entity.Audio) (uint64, error) {
	return s.audioRepo.CreateAudio(audio)
}

func (s *mediaService) GetAudioByID(audioID uint64) (*entity.Audio, string, error) {
	audio, err := s.audioRepo.GetAudioByID(audioID)
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
func (s *mediaService) GetAudioByIDAndUserID(audioID, userID uint64) (*entity.Audio, string, error) {
	// Fetch the audio from the repository
	audio, err := s.audioRepo.GetAudioByIDAndUserID(audioID, userID)
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
func (s *mediaService) ListAudiosByUserID(userID uint64) ([]entity.Audio, error) {
	return s.audioRepo.ListAudiosByUserID(userID)
}

func (s *mediaService) GetAudioByVideoID(videoID, audioID uint64) (*entity.Audio, string, error) {
	audio, err := s.audioRepo.GetAudioByVideoID(videoID, audioID)
	if err != nil {
		return nil, "", err
	}
	presignedURL, err := s.s3Client.GeneratePresignedURL(audio.Folder, audio.FileName, "audio/mpeg")
	if err != nil {
		return nil, "", err
	}
	return audio, presignedURL, nil
}

func (s *mediaService) ListAudiosByVideoID(videoID uint64) ([]entity.Audio, error) {
	return s.audioRepo.ListAudiosByVideoID(videoID)
}

func (s *mediaService) DeleteAudio(audioID uint64) error {
	return s.audioRepo.DeleteAudioByID(audioID)
}

func (s *mediaService) UpdateAudio(audio *entity.Audio) error {
	return s.audioRepo.UpdateAudio(audio)
}

func (s *mediaService) UpdateAudioStatus(audioID uint64, status entity.StatusEntity) error {
	return s.audioRepo.UpdateAudioStatus(audioID, status)
}
