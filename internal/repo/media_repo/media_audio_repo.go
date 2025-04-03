package media_repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
	"strings"
	"time"
)

// CreateAudio inserts a new audio record into the database
func (r *mediaRepo) CreateAudio(audio *entity.Audio) (uint64, error) {
	if audio.Status == "" {
		audio.Status = entity.StatusRaw
	}

	var transcriptionID interface{}
	if audio.TranscriptionID == 0 {
		transcriptionID = nil
	} else {
		transcriptionID = audio.TranscriptionID
	}

	query := `
		INSERT INTO audios (video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		audio.VideoID,
		audio.UserID,
		transcriptionID,
		audio.Duration,
		audio.Lang,
		audio.Folder,
		audio.FileName,
		audio.Status,
		now,
		now,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to execute insert: %v", err)
	}

	// Retrieve the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %v", err)
	}

	return uint64(id), nil
}

// GetAudioByID fetches an audio by its ID
func (r *mediaRepo) GetAudioByID(audioID uint64) (*entity.Audio, error) {
	query := `
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at
		FROM audios
		WHERE id = ?`

	row := r.db.QueryRow(query, audioID)

	var a entity.Audio
	var transcriptionID sql.NullInt64
	err := row.Scan(
		&a.ID,
		&a.VideoID,
		&a.UserID,
		&transcriptionID,
		&a.Duration,
		&a.Lang,
		&a.Folder,
		&a.FileName,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve audio by ID: %v", err)
	}

	if transcriptionID.Valid {
		a.TranscriptionID = uint64(transcriptionID.Int64)
	} else {
		a.TranscriptionID = 0
	}

	return &a, nil
}

// GetAudioByIDAndUserID retrieves a single audio by its ID and User ID (owner)
func (r *mediaRepo) GetAudioByIDAndUserID(audioID, userID uint64) (*entity.Audio, error) {
	query := `
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at
		FROM audios
		WHERE id = ? AND user_id = ?`

	row := r.db.QueryRow(query, audioID, userID)

	var a entity.Audio
	var transcriptionID sql.NullInt64
	err := row.Scan(
		&a.ID,
		&a.VideoID,
		&a.UserID,
		&transcriptionID,
		&a.Duration,
		&a.Lang,
		&a.Folder,
		&a.FileName,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // No record found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve audio by ID and User ID: %v", err)
	}

	if transcriptionID.Valid {
		a.TranscriptionID = uint64(transcriptionID.Int64)
	} else {
		a.TranscriptionID = 0
	}

	return &a, nil
}

// ListAudiosByUserID returns all audios associated with a given user ID
func (r *mediaRepo) ListAudiosByUserID(userID uint64) ([]entity.Audio, error) {
	query := `
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at
		FROM audios
		WHERE user_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audios by user: %v", err)
	}
	defer rows.Close()

	var audios []entity.Audio
	for rows.Next() {
		var a entity.Audio
		var transcriptionID sql.NullInt64
		if err := rows.Scan(
			&a.ID,
			&a.VideoID,
			&a.UserID,
			&transcriptionID,
			&a.Duration,
			&a.Lang,
			&a.Folder,
			&a.FileName,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan audio: %v", err)
		}

		if transcriptionID.Valid {
			a.TranscriptionID = uint64(transcriptionID.Int64)
		} else {
			a.TranscriptionID = 0
		}

		audios = append(audios, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating audio rows: %v", err)
	}

	return audios, nil
}

// GetAudioByVideoID retrieves a specific audio by its video ID and audio ID
func (r *mediaRepo) GetAudioByVideoID(videoID, audioID uint64) (*entity.Audio, error) {
	query := `
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at
		FROM audios
		WHERE video_id = ? AND id = ?`

	row := r.db.QueryRow(query, videoID, audioID)

	var a entity.Audio
	var transcriptionID sql.NullInt64
	err := row.Scan(
		&a.ID,
		&a.VideoID,
		&a.UserID,
		&transcriptionID,
		&a.Duration,
		&a.Lang,
		&a.Folder,
		&a.FileName,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve audio by VideoID and AudioID: %v", err)
	}

	if transcriptionID.Valid {
		a.TranscriptionID = uint64(transcriptionID.Int64)
	} else {
		a.TranscriptionID = 0
	}

	return &a, nil
}

// ListAudiosByVideoID returns all audios associated with a given video ID
func (r *mediaRepo) ListAudiosByVideoID(videoID uint64) ([]entity.Audio, error) {
	query := `
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at
		FROM audios
		WHERE video_id = ?`

	rows, err := r.db.Query(query, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audios by video: %v", err)
	}
	defer rows.Close()

	var audios []entity.Audio
	for rows.Next() {
		var a entity.Audio
		var transcriptionID sql.NullInt64

		if err := rows.Scan(
			&a.ID,
			&a.VideoID,
			&a.UserID,
			&transcriptionID,
			&a.Duration,
			&a.Lang,
			&a.Folder,
			&a.FileName,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan audio: %v", err)
		}

		if transcriptionID.Valid {
			a.TranscriptionID = uint64(transcriptionID.Int64)
		} else {
			a.TranscriptionID = 0
		}

		audios = append(audios, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating audio rows: %v", err)
	}

	return audios, nil
}

// DeleteAudioByID deletes an audio record by its ID
func (r *mediaRepo) DeleteAudioByID(audioID uint64) error {
	query := "DELETE FROM audios WHERE id = ?"
	_, err := r.db.Exec(query, audioID)
	if err != nil {
		return fmt.Errorf("failed to delete audio %d: %v", audioID, err)
	}
	return nil
}

// UpdateAudio updates the entire Audio record
func (r *mediaRepo) UpdateAudio(audio *entity.Audio) error {
	var setClauses []string
	var args []interface{}

	// Dynamically add SET clauses based on non-zero or non-empty fields

	if audio.VideoID != 0 {
		setClauses = append(setClauses, "video_id = ?")
		args = append(args, audio.VideoID)
	}

	if audio.UserID != 0 {
		setClauses = append(setClauses, "user_id = ?")
		args = append(args, audio.UserID)
	}

	if audio.TranscriptionID != 0 {
		setClauses = append(setClauses, "transcription_id = ?")
		args = append(args, audio.TranscriptionID)
	}

	if audio.Duration != 0 {
		setClauses = append(setClauses, "duration = ?")
		args = append(args, audio.Duration)
	}

	if audio.Lang != "" {
		setClauses = append(setClauses, "lang = ?")
		args = append(args, audio.Lang)
	}

	if audio.Folder != "" {
		setClauses = append(setClauses, "folder = ?")
		args = append(args, audio.Folder)
	}

	if audio.FileName != "" {
		setClauses = append(setClauses, "file_name = ?")
		args = append(args, audio.FileName)
	}

	// Always update updated_at
	now := time.Now().Format(time.RFC3339) // Format as RFC3339 string
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, now)

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add the audio.ID for the WHERE clause
	args = append(args, audio.ID)

	// Construct the final SQL query
	query := fmt.Sprintf("UPDATE audios SET %s WHERE id = ?", strings.Join(setClauses, ", "))

	// Execute the query with the arguments
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no audio found with id %d", audio.ID)
	}

	return nil
}

// UpdateAudioStatus updates only the status of an Audio record
func (r *mediaRepo) UpdateAudioStatus(audioID uint64, status entity.StatusEntity) error {
	query := `
        UPDATE audios
        SET status = ?, updated_at = ?
        WHERE id = ?`
	now := time.Now()
	result, err := r.db.Exec(query, status, now, audioID)
	if err != nil {
		return fmt.Errorf("failed to update audio status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no audio found with id %d", audioID)
	}

	return nil
}
