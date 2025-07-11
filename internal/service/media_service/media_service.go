package media_service

import (
	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/pkg/response"
	"mlvt/internal/repo/media_repo"
)

type MediaService interface {
	// audio
	GeneratePresignedUploadURL(folder, fileName, fileType string) (string, error)
	GeneratePresignedDownloadURL(audioID uint64) (string, error)
	CreateAudio(audio *entity.Audio, isFullPipeline bool) (uint64, error)
	GetAudioByID(audioID uint64) (*entity.Audio, string, error)
	GetAudioByIDAndUserID(audioID, userID uint64) (*entity.Audio, string, error)
	ListAudiosByUserID(userID uint64) ([]entity.Audio, error)
	ListAudiosByUserIDAdvance(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity) ([]entity.Audio, error)
	GetAudioByVideoID(videoID, audioID uint64) (*entity.Audio, string, error)
	ListAudiosByVideoID(videoID uint64) ([]entity.Audio, error)
	DeleteAudio(audioID uint64) error
	UpdateAudio(audio *entity.Audio) error
	UpdateAudioStatus(audioID uint64, status entity.StatusEntity) error

	// video
	CreateVideo(video *entity.Video, isFullPipeline bool) (uint64, error)
	GetVideoByID(videoID uint64) (*entity.Video, string, string, error) // Returns the video record and presigned URLs for video and image
	ListVideosByUserID(userID uint64) ([]response.ListVideosByUserIDResponse, error)
	ListVideosByUserIDAdvance(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity) ([]response.ListVideosByUserIDResponse, error)
	DeleteVideo(videoID uint64) error
	UpdateVideo(video *entity.Video) error
	UpdateVideoStatus(videoID uint64, status entity.StatusEntity) error
	GetVideoStatus(videoID uint64) (entity.StatusEntity, error)
	GeneratePresignedUploadURLForVideo(folder, fileName, fileType string) (string, error)
	GeneratePresignedUploadURLForImage(folder, fileName, fileType string) (string, error)
	GeneratePresignedDownloadURLForVideo(videoID uint64) (string, error)
	GeneratePresignedDownloadURLForImage(videoID uint64) (string, error)

	// transcription
	CreateTranscription(transcription *entity.Transcription, isFullPipeline bool, isOriginalText bool) (uint64, error)
	GetTranscriptionByID(transcriptionID uint64) (*entity.Transcription, string, error)
	GetTranscriptionByIDAndUserID(transcriptionID, userID uint64) (*entity.Transcription, string, error)
	GetTranscriptionByIDAndVideoID(transcriptionID, videoID uint64) (*entity.Transcription, string, error)
	ListTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error)
	ListTranscriptionsByUserIDAdvance(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity) ([]entity.Transcription, error)
	ListTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error)
	DeleteTranscription(transcriptionID uint64) error
	GeneratePresignedUploadURLForText(folder, fileName, fileType string) (string, error)
	GeneratePresignedDownloadURLForText(transcriptionID uint64) (string, error)
	UpdateTranscription(transcription *entity.Transcription) error
	UpdateTranscriptionStatus(transcriptionID uint64, status entity.StatusEntity) error

	// GetAll media
	GetAllMedia(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity, mediaType []entity.MediaType) (response.ProcessResponse, error)
}

type mediaService struct {
	mediaRepo media_repo.MediaRepository
	s3Client  aws.S3ClientInterface
}

func NewMediaService(
	mediaRepo media_repo.MediaRepository,
	s3Client aws.S3ClientInterface,
) MediaService {
	return &mediaService{
		mediaRepo: mediaRepo,
		s3Client:  s3Client,
	}
}
