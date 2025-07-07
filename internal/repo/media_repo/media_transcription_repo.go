package media_repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// CreateTranscription inserts a new transcription into the database
func (r *mediaRepo) CreateTranscription(transcription *entity.Transcription) (uint64, error) {
	if transcription.Status == "" {
		transcription.Status = entity.StatusRaw
	}

	var videoID interface{}
	if transcription.VideoID == 0 {
		videoID = nil
	} else {
		videoID = transcription.VideoID
	}

	var originalTranscriptionID interface{}
	if transcription.OriginalTranscriptionID == 0 {
		originalTranscriptionID = nil
	} else {
		originalTranscriptionID = transcription.OriginalTranscriptionID
	}

	query := `
		INSERT INTO transcriptions (
			video_id, user_id, original_transcription_id, text,
			lang, folder, file_name, status, created_at, updated_at, title
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	now := time.Now()
	var insertedID uint64
	err := r.db.QueryRow(
		query,
		videoID,
		transcription.UserID,
		originalTranscriptionID,
		transcription.Text,
		transcription.Lang,
		transcription.Folder,
		transcription.FileName,
		transcription.Status,
		now,
		now,
		transcription.Title,
	).Scan(&insertedID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert transcription: %w", err)
	}

	return insertedID, nil
}

// GetTranscriptionByID retrieves a transcription by its ID
func (r *mediaRepo) GetTranscriptionByID(transcriptionID uint64) (*entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at, title
		FROM transcriptions
		WHERE id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	row := r.db.QueryRow(query, transcriptionID)

	var t entity.Transcription
	var videoID sql.NullInt64
	var originalTranscriptionID sql.NullInt64

	err := row.Scan(
		&t.ID,
		&videoID,
		&t.UserID,
		&originalTranscriptionID,
		&t.Text,
		&t.Lang,
		&t.Folder,
		&t.FileName,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.Title,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transcription by ID: %w", err)
	}

	if videoID.Valid {
		t.VideoID = uint64(videoID.Int64)
	} else {
		t.VideoID = 0
	}
	if originalTranscriptionID.Valid {
		t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
	} else {
		t.OriginalTranscriptionID = 0
	}

	return &t, nil
}


// GetTranscriptionByIDAndUserID retrieves a transcription by its ID and User ID
func (r *mediaRepo) GetTranscriptionByIDAndUserID(transcriptionID, userID uint64) (*entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at, title
		FROM transcriptions
		WHERE id = ? AND user_id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	row := r.db.QueryRow(query, transcriptionID, userID)

	var t entity.Transcription
	var videoID sql.NullInt64
	var originalTranscriptionID sql.NullInt64

	err := row.Scan(
		&t.ID,
		&videoID,
		&t.UserID,
		&originalTranscriptionID,
		&t.Text,
		&t.Lang,
		&t.Folder,
		&t.FileName,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.Title,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transcription by ID and User ID: %w", err)
	}

	if videoID.Valid {
		t.VideoID = uint64(videoID.Int64)
	} else {
		t.VideoID = 0
	}
	if originalTranscriptionID.Valid {
		t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
	} else {
		t.OriginalTranscriptionID = 0
	}

	return &t, nil
}

// GetTranscriptionByIDAndVideoID retrieves a transcription by its ID and Video ID
func (r *mediaRepo) GetTranscriptionByIDAndVideoID(transcriptionID, videoID uint64) (*entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at, title
		FROM transcriptions
		WHERE id = ? AND video_id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	row := r.db.QueryRow(query, transcriptionID, videoID)

	var t entity.Transcription
	var videoIDNullable sql.NullInt64
	var originalTranscriptionID sql.NullInt64

	err := row.Scan(
		&t.ID,
		&videoIDNullable,
		&t.UserID,
		&originalTranscriptionID,
		&t.Text,
		&t.Lang,
		&t.Folder,
		&t.FileName,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.Title,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transcription by ID and Video ID: %w", err)
	}

	if videoIDNullable.Valid {
		t.VideoID = uint64(videoIDNullable.Int64)
	} else {
		t.VideoID = 0
	}
	if originalTranscriptionID.Valid {
		t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
	} else {
		t.OriginalTranscriptionID = 0
	}

	return &t, nil
}

// ListTranscriptionsByUserID lists all transcriptions for a specific user
func (r *mediaRepo) ListTranscriptionsByUserID(userID uint64) ([]entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at, title
		FROM transcriptions
		WHERE user_id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transcriptions by user: %w", err)
	}
	defer rows.Close()

	var transcriptions []entity.Transcription
	for rows.Next() {
		var t entity.Transcription
		var videoID sql.NullInt64
		var originalTranscriptionID sql.NullInt64

		if err := rows.Scan(
			&t.ID,
			&videoID,
			&t.UserID,
			&originalTranscriptionID,
			&t.Text,
			&t.Lang,
			&t.Folder,
			&t.FileName,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.Title,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transcription: %w", err)
		}

		if videoID.Valid {
			t.VideoID = uint64(videoID.Int64)
		} else {
			t.VideoID = 0
		}
		if originalTranscriptionID.Valid {
			t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
		} else {
			t.OriginalTranscriptionID = 0
		}

		transcriptions = append(transcriptions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transcription rows: %w", err)
	}

	return transcriptions, nil
}

// ListTranscriptionsByUserID lists all transcriptions for a specific user
func (r *mediaRepo) ListTranscriptionsByUserIDAdvance(
	userID uint64,
	searchKey string,
	limit int,
	offset int,
	status []entity.StatusEntity,
) ([]entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at, title
		FROM transcriptions
		WHERE user_id = ? AND title ILIKE ? AND is_deleted = FALSE AND status IN (?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	searchKey = "%" + searchKey + "%"
	query, args, err := sqlx.In(query, userID, searchKey, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to build IN clause: %w", err)
	}
	query = r.db.Rebind(query)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query transcriptions by user: %w", err)
	}
	defer rows.Close()

	var transcriptions []entity.Transcription
	for rows.Next() {
		var t entity.Transcription
		var videoID sql.NullInt64
		var originalTranscriptionID sql.NullInt64

		if err := rows.Scan(
			&t.ID,
			&videoID,
			&t.UserID,
			&originalTranscriptionID,
			&t.Text,
			&t.Lang,
			&t.Folder,
			&t.FileName,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.Title,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transcription: %w", err)
		}

		if videoID.Valid {
			t.VideoID = uint64(videoID.Int64)
		} else {
			t.VideoID = 0
		}
		if originalTranscriptionID.Valid {
			t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
		} else {
			t.OriginalTranscriptionID = 0
		}

		transcriptions = append(transcriptions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transcription rows: %w", err)
	}

	return transcriptions, nil
}

// ListTranscriptionsByVideoID lists all transcriptions for a specific video
func (r *mediaRepo) ListTranscriptionsByVideoID(videoID uint64) ([]entity.Transcription, error) {
	query := `
		SELECT id, video_id, user_id, original_transcription_id, text, lang, folder, file_name, status, created_at, updated_at, title
		FROM transcriptions
		WHERE video_id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := r.db.Query(query, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transcriptions by video: %w", err)
	}
	defer rows.Close()

	var transcriptions []entity.Transcription
	for rows.Next() {
		var t entity.Transcription
		var videoIDNullable sql.NullInt64
		var originalTranscriptionID sql.NullInt64

		if err := rows.Scan(
			&t.ID,
			&videoIDNullable,
			&t.UserID,
			&originalTranscriptionID,
			&t.Text,
			&t.Lang,
			&t.Folder,
			&t.FileName,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.Title,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transcription: %w", err)
		}

		if videoIDNullable.Valid {
			t.VideoID = uint64(videoIDNullable.Int64)
		} else {
			t.VideoID = 0
		}
		if originalTranscriptionID.Valid {
			t.OriginalTranscriptionID = uint64(originalTranscriptionID.Int64)
		} else {
			t.OriginalTranscriptionID = 0
		}

		transcriptions = append(transcriptions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transcription rows: %w", err)
	}

	return transcriptions, nil
}

// DeleteTranscription deletes a transcription by its ID
func (r *mediaRepo) DeleteTranscription(transcriptionID uint64) error {
	query := "UPDATE transcriptions SET is_deleted = true WHERE id = ?"
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	_, err := r.db.Exec(query, transcriptionID)
	if err != nil {
		return fmt.Errorf("failed to delete transcription %d: %w", transcriptionID, err)
	}
	return nil
}

func (r *mediaRepo) GetCountTranscriptionsByUserId(userID uint64) (int, error) {
	query := `SELECT count(*) FROM transcriptions WHERE user_id = ?`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	row := r.db.QueryRow(query, userID)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to scan count: %w", err)
	}

	return count, nil
}

// UpdateTranscription updates an existing transcription record
func (r *mediaRepo) UpdateTranscription(transcription *entity.Transcription) error {

	var setClauses []string
	var args []interface{}

	if transcription.OriginalTranscriptionID != 0 {
		setClauses = append(setClauses, "original_transcription_id = ?")
		args = append(args, transcription.OriginalTranscriptionID)
	}
	if transcription.Text != "" {
		setClauses = append(setClauses, "text = ?")
		args = append(args, transcription.Text)
	}
	if transcription.Title != "" {
		setClauses = append(setClauses, "title = ?")
		args = append(args, transcription.Title)
	}
	if transcription.Lang != "" {
		setClauses = append(setClauses, "lang = ?")
		args = append(args, transcription.Lang)
	}

	// Always update updated_at
	now := time.Now()
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, now)

	query := fmt.Sprintf("UPDATE transcriptions SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no transcription found with id %d", transcription.ID)
	}
	return nil
}

// UpdateTranscriptionStatus updates only the status of a transcription record
func (r *mediaRepo) UpdateTranscriptionStatus(transcriptionID uint64, status entity.StatusEntity) error {
	query := `
		UPDATE transcriptions
		SET status = ?, updated_at = ?
		WHERE id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	now := time.Now()
	result, err := r.db.Exec(query, status, now, transcriptionID)
	if err != nil {
		return fmt.Errorf("failed to update transcription status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no transcription found with id %d", transcriptionID)
	}

	return nil
}
