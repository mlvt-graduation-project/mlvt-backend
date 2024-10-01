package repo

import (
	"database/sql"
	"mlvt/internal/entity"
	"time"
)

type TranscriptionRepository interface {
	CreateTranscription(transcription *entity.Transcription) error
	GetTranscriptionByID(transcriptionID uint64) (*entity.Transcription, error)
	GetTranscriptionByIDAndUserID(transcriptionID, userID uint64) (*entity.Transcription, error)
	GetTranscriptionByIDAndVideoID(transcriptionID, videoID uint64) (*entity.Transcription, error)
	ListTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error)
	ListTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error)
	DeleteTranscription(transcriptionID uint64) error
}

type transcriptionRepo struct {
	db *sql.DB
}

func NewTranscriptionRepository(db *sql.DB) TranscriptionRepository {
	return &transcriptionRepo{db: db}
}

// CreateTranscription inserts a new transcription into the database
func (r *transcriptionRepo) CreateTranscription(transcription *entity.Transcription) error {
	query := `
		INSERT INTO transcriptions (video_id, user_id, text, lang, folder, file_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	_, err := r.db.Exec(query, transcription.VideoID, transcription.UserID, transcription.Text,
		transcription.Lang, transcription.Folder, transcription.FileName, now, now)
	return err
}

// GetTranscriptionByID retrieves a transcription by its ID
func (r *transcriptionRepo) GetTranscriptionByID(transcriptionID uint64) (*entity.Transcription, error) {
	query := `SELECT id, video_id, user_id, text, lang, folder, file_name, created_at, updated_at
	          FROM transcriptions WHERE id = ?`
	row := r.db.QueryRow(query, transcriptionID)
	transcription := &entity.Transcription{}
	err := row.Scan(&transcription.ID, &transcription.VideoID, &transcription.UserID, &transcription.Text,
		&transcription.Lang, &transcription.Folder, &transcription.FileName, &transcription.CreatedAt, &transcription.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return transcription, err
}

// GetTranscriptionByIDAndUserID retrieves a transcription by its ID and User ID
func (r *transcriptionRepo) GetTranscriptionByIDAndUserID(transcriptionID, userID uint64) (*entity.Transcription, error) {
	query := `SELECT id, video_id, user_id, text, lang, folder, file_name, created_at, updated_at
	          FROM transcriptions WHERE id = ? AND user_id = ?`
	row := r.db.QueryRow(query, transcriptionID, userID)
	transcription := &entity.Transcription{}
	err := row.Scan(&transcription.ID, &transcription.VideoID, &transcription.UserID, &transcription.Text,
		&transcription.Lang, &transcription.Folder, &transcription.FileName, &transcription.CreatedAt, &transcription.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return transcription, err
}

// GetTranscriptionByIDAndVideoID retrieves a transcription by its ID and Video ID
func (r *transcriptionRepo) GetTranscriptionByIDAndVideoID(transcriptionID, videoID uint64) (*entity.Transcription, error) {
	query := `SELECT id, video_id, user_id, text, lang, folder, file_name, created_at, updated_at
	          FROM transcriptions WHERE id = ? AND video_id = ?`
	row := r.db.QueryRow(query, transcriptionID, videoID)
	transcription := &entity.Transcription{}
	err := row.Scan(&transcription.ID, &transcription.VideoID, &transcription.UserID, &transcription.Text,
		&transcription.Lang, &transcription.Folder, &transcription.FileName, &transcription.CreatedAt, &transcription.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return transcription, err
}

// ListTranscriptionsByUserID lists all transcriptions for a specific user
func (r *transcriptionRepo) ListTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error) {
	query := `SELECT id, video_id, user_id, text, lang, folder, file_name, created_at, updated_at
	          FROM transcriptions WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transcriptions []entity.Transcription
	for rows.Next() {
		var transcription entity.Transcription
		if err := rows.Scan(&transcription.ID, &transcription.VideoID, &transcription.UserID, &transcription.Text,
			&transcription.Lang, &transcription.Folder, &transcription.FileName, &transcription.CreatedAt, &transcription.UpdatedAt); err != nil {
			return nil, err
		}
		transcriptions = append(transcriptions, transcription)
	}
	return transcriptions, nil
}

// ListTranscriptionsByVideoID lists all transcriptions for a specific video
func (r *transcriptionRepo) ListTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error) {
	query := `SELECT id, video_id, user_id, text, lang, folder, file_name, created_at, updated_at
	          FROM transcriptions WHERE video_id = ?`
	rows, err := r.db.Query(query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transcriptions []entity.Transcription
	for rows.Next() {
		var transcription entity.Transcription
		if err := rows.Scan(&transcription.ID, &transcription.VideoID, &transcription.UserID, &transcription.Text,
			&transcription.Lang, &transcription.Folder, &transcription.FileName, &transcription.CreatedAt, &transcription.UpdatedAt); err != nil {
			return nil, err
		}
		transcriptions = append(transcriptions, transcription)
	}
	return transcriptions, nil
}

// DeleteTranscription deletes a transcription by its ID
func (r *transcriptionRepo) DeleteTranscription(transcriptionID uint64) error {
	query := "DELETE FROM transcriptions WHERE id = ?"
	_, err := r.db.Exec(query, transcriptionID)
	return err
}
