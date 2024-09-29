package repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
)

type TranscriptionRepository interface {
	CreateTranscription(tx entity.Transcription) (entity.Transcription, error)
	GetTranscriptionByID(userID, transcriptionID uint64) (entity.Transcription, error)
	GetTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error)
	GetTranscriptionByVideoID(videoID, transcriptionID uint64) (entity.Transcription, error)
	GetTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error)
	DeleteTranscription(transcriptionID uint64) error
}

type transcriptionRepo struct {
	DB *sql.DB
}

func NewTranscriptionRepository(db *sql.DB) *transcriptionRepo {
	return &transcriptionRepo{DB: db}
}

func (r *transcriptionRepo) CreateTranscription(tx entity.Transcription) (entity.Transcription, error) {
	query := `INSERT INTO transcriptions (video_id, user_id, text, lang, folder, file_name, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW()) RETURNING id`
	err := r.DB.QueryRow(query, tx.VideoID, tx.UserID, tx.Text, tx.Lang, tx.Folder, tx.FileName).Scan(&tx.ID)
	if err != nil {
		return entity.Transcription{}, fmt.Errorf("error inserting transcription: %w", err)
	}
	return tx, nil
}

func (r *transcriptionRepo) GetTranscriptionByID(userID, transcriptionID uint64) (entity.Transcription, error) {
	query := `SELECT id, video_id, user_id, text, lang, folder, file_name, created_at, updated_at
			  FROM transcriptions WHERE id = ? AND user_id = ?`
	var tx entity.Transcription
	err := r.DB.QueryRow(query, transcriptionID, userID).Scan(&tx.ID, &tx.VideoID, &tx.UserID, &tx.Text, &tx.Lang, &tx.Folder, &tx.FileName, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		return entity.Transcription{}, fmt.Errorf("error retrieving transcription: %w", err)
	}
	return tx, nil
}

func (r *transcriptionRepo) GetTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error) {
	query := `SELECT id, video_id, user_id, text, lang, folder, file_name, created_at, updated_at
			  FROM transcriptions WHERE user_id = ?`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transcriptions by user ID: %w", err)
	}
	defer rows.Close()

	var transcriptions []entity.Transcription
	for rows.Next() {
		var tx entity.Transcription
		if err := rows.Scan(&tx.ID, &tx.VideoID, &tx.UserID, &tx.Text, &tx.Lang, &tx.Folder, &tx.FileName, &tx.CreatedAt, &tx.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning transcription: %w", err)
		}
		transcriptions = append(transcriptions, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error finalizing transcription retrieval: %w", err)
	}
	return transcriptions, nil
}

func (r *transcriptionRepo) GetTranscriptionByVideoID(videoID, transcriptionID uint64) (entity.Transcription, error) {
	query := `SELECT * FROM transcriptions WHERE video_id = ? AND id = ?`
	var tx entity.Transcription
	err := r.DB.QueryRow(query, videoID, transcriptionID).Scan(&tx.ID, &tx.VideoID, &tx.UserID, &tx.Text, &tx.Lang, &tx.Folder, &tx.FileName, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		return entity.Transcription{}, fmt.Errorf("error retrieving transcription by video ID: %w", err)
	}
	return tx, nil
}

func (r *transcriptionRepo) GetTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error) {
	query := `SELECT * FROM transcriptions WHERE video_id = ?`
	rows, err := r.DB.Query(query, videoID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transcriptions by video ID: %w", err)
	}
	defer rows.Close()

	var transcriptions []entity.Transcription
	for rows.Next() {
		var tx entity.Transcription
		if err := rows.Scan(&tx.ID, &tx.VideoID, &tx.UserID, &tx.Text, &tx.Lang, &tx.Folder, &tx.FileName, &tx.CreatedAt, &tx.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning transcription: %w", err)
		}
		transcriptions = append(transcriptions, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("finalizing retrieval of transcriptions: %w", err)
	}
	return transcriptions, nil
}

func (r *transcriptionRepo) DeleteTranscription(transcriptionID uint64) error {
	query := `DELETE FROM transcriptions WHERE id = ?`
	_, err := r.DB.Exec(query, transcriptionID)
	if err != nil {
		return fmt.Errorf("error deleting transcription: %w", err)
	}
	return nil
}
