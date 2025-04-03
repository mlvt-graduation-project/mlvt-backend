package media_repo

import (
	"database/sql"
	"mlvt/internal/entity"
)

type MediaRepository interface {
	// audio
	CreateAudio(audio *entity.Audio) (uint64, error)
	GetAudioByID(audioID uint64) (*entity.Audio, error)
	GetAudioByIDAndUserID(audioID, userID uint64) (*entity.Audio, error)
	ListAudiosByUserID(userID uint64) ([]entity.Audio, error)
	GetAudioByVideoID(videoID, audioID uint64) (*entity.Audio, error)
	ListAudiosByVideoID(videoID uint64) ([]entity.Audio, error)
	DeleteAudioByID(audioID uint64) error
	UpdateAudio(audio *entity.Audio) error
	UpdateAudioStatus(audioID uint64, status entity.StatusEntity) error

	// video
	CreateVideo(video *entity.Video) (uint64, error)
	GetVideoByID(videoID uint64) (*entity.Video, error)
	ListVideosByUserID(userID uint64) ([]entity.Video, error)
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
	ListTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error)
	DeleteTranscription(transcriptionID uint64) error
	UpdateTranscription(transcription *entity.Transcription) error
	UpdateTranscriptionStatus(transcriptionID uint64, status entity.StatusEntity) error
}

type mediaRepo struct {
	db *sql.DB
}

func NewMediaRepo(db *sql.DB) MediaRepository {
	return &mediaRepo{db: db}
}
