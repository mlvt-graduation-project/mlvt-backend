package media_service

import (
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/utility"
)

func (s *mediaService) GeneratePresignedUploadURL(folder, fileName, fileType string) (string, error) {
	return s.s3Client.GeneratePresignedURL(folder, fileName, fileType)
}

func (s *mediaService) GeneratePresignedDownloadURL(audioID uint64) (string, error) {
	// Fetch the audio from the repository using its ID
	audio, err := s.mediaRepo.GetAudioByID(audioID)
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

func (s *mediaService) CreateAudio(audio *entity.Audio, isFullPipeline bool) (uint64, error) {
	count, err := s.mediaRepo.GetCountAudiosByUserId(audio.UserID)
	if err != nil {
		return 0, fmt.Errorf("error querying total audios of user id")
	}
	audio.Title = utility.GetMediaTitle(entity.MediaTypeAudio, isFullPipeline, false, count+1)

	return s.mediaRepo.CreateAudio(audio)
}

func (s *mediaService) GetAudioByID(audioID uint64) (*entity.Audio, string, error) {
	audio, err := s.mediaRepo.GetAudioByID(audioID)
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
	audio, err := s.mediaRepo.GetAudioByIDAndUserID(audioID, userID)
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
	return s.mediaRepo.ListAudiosByUserID(userID)
}

func (s *mediaService) ListAudiosByUserIDAdvance(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity) ([]entity.Audio, error) {
	return s.mediaRepo.ListAudiosByUserIDAdvance(userID, searchKey, limit, offset, status)
}

func (s *mediaService) GetAudioByVideoID(videoID, audioID uint64) (*entity.Audio, string, error) {
	audio, err := s.mediaRepo.GetAudioByVideoID(videoID, audioID)
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
	return s.mediaRepo.ListAudiosByVideoID(videoID)
}

func (s *mediaService) DeleteAudio(audioID uint64) error {
	return s.mediaRepo.DeleteAudioByID(audioID)
}

func (s *mediaService) UpdateAudio(audio *entity.Audio) error {
	return s.mediaRepo.UpdateAudio(audio)
}

func (s *mediaService) UpdateAudioStatus(audioID uint64, status entity.StatusEntity) error {
	return s.mediaRepo.UpdateAudioStatus(audioID, status)
}
