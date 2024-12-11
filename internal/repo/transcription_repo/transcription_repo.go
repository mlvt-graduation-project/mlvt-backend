package transcription_repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
	"time"
)

type TranscriptionRepository interface {
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

type transcriptionRepo struct {
	db *sql.DB
}

func NewTranscriptionRepository(db *sql.DB) TranscriptionRepository {
	return &transcriptionRepo{db: db}
}

// CreateTranscription inserts a new transcription into the database
func (r *transcriptionRepo) CreateTranscription(transcription *entity.Transcription) (uint64, error) {
	if transcription.Status == "" {
		transcription.Status = entity.StatusRaw
	}

	var originalTranscriptionID interface{}
	if transcription.OriginalTranscriptionID == 0 {
		originalTranscriptionID = nil
	} else {
		originalTranscriptionID = transcription.OriginalTranscriptionID
	}

	query := `
		INSERT INTO transcriptions (video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := r.db.Exec(query,
		transcription.VideoID,
		transcription.UserID,
		originalTranscriptionID,
		transcription.Text,
		transcription.Lang,
		transcription.Folder,
		transcription.FileName,
		transcription.Status,
		now,
		now,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to execute insert: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %v", err)
	}

	return uint64(id), nil
}

// GetTranscriptionByID retrieves a transcription by its ID
func (r *transcriptionRepo) GetTranscriptionByID(transcriptionID uint64) (*entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at
		FROM transcriptions
		WHERE id = ?`
	row := r.db.QueryRow(query, transcriptionID)

	var t entity.Transcription
	var originalTranscriptionID sql.NullInt64

	err := row.Scan(
		&t.ID,
		&t.VideoID,
		&t.UserID,
		&originalTranscriptionID,
		&t.Text,
		&t.Lang,
		&t.Folder,
		&t.FileName,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transcription by ID: %v", err)
	}

	if originalTranscriptionID.Valid {
		t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
	} else {
		t.OriginalTranscriptionID = 0
	}

	return &t, nil
}

// GetTranscriptionByIDAndUserID retrieves a transcription by its ID and User ID
func (r *transcriptionRepo) GetTranscriptionByIDAndUserID(transcriptionID, userID uint64) (*entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at
		FROM transcriptions
		WHERE id = ? AND user_id = ?`
	row := r.db.QueryRow(query, transcriptionID, userID)

	var t entity.Transcription
	var originalTranscriptionID sql.NullInt64

	err := row.Scan(
		&t.ID,
		&t.VideoID,
		&t.UserID,
		&originalTranscriptionID,
		&t.Text,
		&t.Lang,
		&t.Folder,
		&t.FileName,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transcription by ID and User ID: %v", err)
	}

	if originalTranscriptionID.Valid {
		t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
	} else {
		t.OriginalTranscriptionID = 0
	}

	return &t, nil
}

// GetTranscriptionByIDAndVideoID retrieves a transcription by its ID and Video ID
func (r *transcriptionRepo) GetTranscriptionByIDAndVideoID(transcriptionID, videoID uint64) (*entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at
		FROM transcriptions
		WHERE id = ? AND video_id = ?`
	row := r.db.QueryRow(query, transcriptionID, videoID)

	var t entity.Transcription
	var originalTranscriptionID sql.NullInt64

	err := row.Scan(
		&t.ID,
		&t.VideoID,
		&t.UserID,
		&originalTranscriptionID,
		&t.Text,
		&t.Lang,
		&t.Folder,
		&t.FileName,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transcription by ID and Video ID: %v", err)
	}

	if originalTranscriptionID.Valid {
		t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
	} else {
		t.OriginalTranscriptionID = 0
	}

	return &t, nil
}

// ListTranscriptionsByUserID lists all transcriptions for a specific user
func (r *transcriptionRepo) ListTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at
		FROM transcriptions
		WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transcriptions by user: %v", err)
	}
	defer rows.Close()

	var transcriptions []entity.Transcription
	for rows.Next() {
		var t entity.Transcription
		var originalTranscriptionID sql.NullInt64

		if err := rows.Scan(
			&t.ID,
			&t.VideoID,
			&t.UserID,
			&originalTranscriptionID,
			&t.Text,
			&t.Lang,
			&t.Folder,
			&t.FileName,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transcription: %v", err)
		}

		if originalTranscriptionID.Valid {
			t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
		} else {
			t.OriginalTranscriptionID = 0
		}

		transcriptions = append(transcriptions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transcription rows: %v", err)
	}

	return transcriptions, nil
}

// ListTranscriptionsByVideoID lists all transcriptions for a specific video
func (r *transcriptionRepo) ListTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at
		FROM transcriptions
		WHERE video_id = ?`
	rows, err := r.db.Query(query, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transcriptions by video: %v", err)
	}
	defer rows.Close()

	var transcriptions []entity.Transcription
	for rows.Next() {
		var t entity.Transcription
		var originalTranscriptionID sql.NullInt64

		if err := rows.Scan(
			&t.ID,
			&t.VideoID,
			&t.UserID,
			&originalTranscriptionID,
			&t.Text,
			&t.Lang,
			&t.Folder,
			&t.FileName,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transcription: %v", err)
		}

		if originalTranscriptionID.Valid {
			t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
		} else {
			t.OriginalTranscriptionID = 0
		}

		transcriptions = append(transcriptions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transcription rows: %v", err)
	}

	return transcriptions, nil
}

// DeleteTranscription deletes a transcription by its ID
func (r *transcriptionRepo) DeleteTranscription(transcriptionID uint64) error {
	query := "DELETE FROM transcriptions WHERE id = ?"
	_, err := r.db.Exec(query, transcriptionID)
	if err != nil {
		return fmt.Errorf("failed to delete transcription %d: %v", transcriptionID, err)
	}
	return nil
}

// UpdateTranscription updates an existing transcription record
func (r *transcriptionRepo) UpdateTranscription(transcription *entity.Transcription) error {
	var originalTranscriptionID interface{}
	if transcription.OriginalTranscriptionID == 0 {
		originalTranscriptionID = nil
	} else {
		originalTranscriptionID = transcription.OriginalTranscriptionID
	}

	query := `
		UPDATE transcriptions
		SET original_transcription_id = ?, text = ?, lang = ?, folder = ?, file_name = ?, updated_at = ?
		WHERE id = ?`
	now := time.Now()
	result, err := r.db.Exec(query,
		originalTranscriptionID,
		transcription.Text,
		transcription.Lang,
		transcription.Folder,
		transcription.FileName,
		now,
		transcription.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to execute update: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no transcription found with id %d", transcription.ID)
	}
	return nil
}

// UpdateTranscriptionStatus updates only the status of a transcription record
func (r *transcriptionRepo) UpdateTranscriptionStatus(transcriptionID uint64, status entity.StatusEntity) error {
	query := `
		UPDATE transcriptions
		SET status = ?, updated_at = ?
		WHERE id = ?`
	now := time.Now()
	result, err := r.db.Exec(query, status, now, transcriptionID)
	if err != nil {
		return fmt.Errorf("failed to update transcription status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no transcription found with id %d", transcriptionID)
	}

	return nil
}
