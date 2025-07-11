package media_repo

import (
	"mlvt/internal/entity"

	"github.com/jmoiron/sqlx"
)

type MediaRepository interface {
	// audio
	CreateAudio(audio *entity.Audio) (uint64, error)
	GetAudioByID(audioID uint64) (*entity.Audio, error)
	GetAudioByIDAndUserID(audioID, userID uint64) (*entity.Audio, error)
	ListAudiosByUserID(userID uint64) ([]entity.Audio, error)
	ListAudiosByUserIDAdvance(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity) ([]entity.Audio, error)
	GetCountAudiosByUserId(userID uint64) (int, error)
	GetAudioByVideoID(videoID, audioID uint64) (*entity.Audio, error)
	ListAudiosByVideoID(videoID uint64) ([]entity.Audio, error)
	DeleteAudioByID(audioID uint64) error
	UpdateAudio(audio *entity.Audio) error
	UpdateAudioStatus(audioID uint64, status entity.StatusEntity) error

	// video
	CreateVideo(video *entity.Video) (uint64, error)
	GetVideoByID(videoID uint64) (*entity.Video, error)
	ListVideosByUserID(userID uint64) ([]entity.Video, error)
	ListVideosByUserIDAdvance(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity) ([]entity.Video, error)
	GetCountVideosByUserId(userID uint64) (int, error)
	DeleteVideo(videoID uint64) error
	UpdateVideo(video *entity.Video) error
	GetVideoStatus(videoID uint64) (entity.StatusEntity, error)
	UpdateVideoStatus(videoId uint64, status entity.StatusEntity) error

	// transcription
	CreateTranscription(transcription *entity.Transcription) (uint64, error)
	GetTranscriptionByID(transcriptionID uint64) (*entity.Transcription, error)
	GetTranscriptionByIDAndUserID(transcriptionID, userID uint64) (*entity.Transcription, error)
	GetTranscriptionByIDAndVideoID(transcriptionID, videoID uint64) (*entity.Transcription, error)
	ListTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error)
	ListTranscriptionsByUserIDAdvance(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity) ([]entity.Transcription, error)
	GetCountTranscriptionsByUserId(userID uint64) (int, error)
	ListTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error)
	DeleteTranscription(transcriptionID uint64) error
	UpdateTranscription(transcription *entity.Transcription) error
	UpdateTranscriptionStatus(transcriptionID uint64, status entity.StatusEntity) error

	// media
	GetAllMedia(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity, mediaType []entity.MediaType) ([]entity.Video, []entity.Audio, []entity.Transcription, int, error)
}

type mediaRepo struct {
	db *sqlx.DB
}

func NewMediaRepo(db *sqlx.DB) MediaRepository {
	return &mediaRepo{db: db}
}
