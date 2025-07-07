package media_repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
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

	var videoID interface{}
	if audio.VideoID == 0 {
		videoID = nil
	} else {
		videoID = audio.VideoID
	}

	now := time.Now()

	query := `
		INSERT INTO audios (
			video_id, user_id, transcription_id, duration,
			lang, folder, file_name, status, created_at, updated_at, title
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`

	// Rebind to $1, $2... for PostgreSQL
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	var insertedID uint64
	err := r.db.QueryRow(query,
		videoID,
		audio.UserID,
		transcriptionID,
		audio.Duration,
		audio.Lang,
		audio.Folder,
		audio.FileName,
		audio.Status,
		now,
		now,
		audio.Title,
	).Scan(&insertedID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert audio: %w", err)
	}

	return insertedID, nil
}


// GetAudioByID fetches an audio by its ID
func (r *mediaRepo) GetAudioByID(audioID uint64) (*entity.Audio, error) {
	query := `
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at, title
		FROM audios
		WHERE id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query) // convert ? → $1

	row := r.db.QueryRow(query, audioID)

	var a entity.Audio
	var transcriptionID sql.NullInt64
	var videoID sql.NullInt64

	err := row.Scan(
		&a.ID,
		&videoID,
		&a.UserID,
		&transcriptionID,
		&a.Duration,
		&a.Lang,
		&a.Folder,
		&a.FileName,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Title,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve audio by ID: %w", err)
	}

	if transcriptionID.Valid {
		a.TranscriptionID = uint64(transcriptionID.Int64)
	} else {
		a.TranscriptionID = 0
	}

	if videoID.Valid {
		a.VideoID = uint64(videoID.Int64)
	} else {
		a.TranscriptionID = 0
	}

	return &a, nil
}


// GetAudioByIDAndUserID retrieves a single audio by its ID and User ID (owner)
func (r *mediaRepo) GetAudioByIDAndUserID(audioID, userID uint64) (*entity.Audio, error) {
	query := `
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at, title
		FROM audios
		WHERE id = ? AND user_id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query) // ? → $1, $2

	row := r.db.QueryRow(query, audioID, userID)

	var a entity.Audio
	var videoID sql.NullInt64
	var transcriptionID sql.NullInt64

	err := row.Scan(
		&a.ID,
		&videoID,
		&a.UserID,
		&transcriptionID,
		&a.Duration,
		&a.Lang,
		&a.Folder,
		&a.FileName,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Title,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve audio by ID and User ID: %w", err)
	}

	// Handle possible NULLs
	if videoID.Valid {
		a.VideoID = uint64(videoID.Int64)
	} else {
		a.VideoID = 0
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
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at, title
		FROM audios
		WHERE user_id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audios by user: %w", err)
	}
	defer rows.Close()

	var audios []entity.Audio
	for rows.Next() {
		var a entity.Audio
		var videoID sql.NullInt64
		var transcriptionID sql.NullInt64

		if err := rows.Scan(
			&a.ID,
			&videoID,
			&a.UserID,
			&transcriptionID,
			&a.Duration,
			&a.Lang,
			&a.Folder,
			&a.FileName,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.Title,
		); err != nil {
			return nil, fmt.Errorf("failed to scan audio: %w", err)
		}

		if videoID.Valid {
			a.VideoID = uint64(videoID.Int64)
		} else {
			a.VideoID = 0
		}

		if transcriptionID.Valid {
			a.TranscriptionID = uint64(transcriptionID.Int64)
		} else {
			a.TranscriptionID = 0
		}

		audios = append(audios, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating audio rows: %w", err)
	}

	return audios, nil
}

func (r *mediaRepo) GetCountAudiosByUserId(userID uint64) (int, error) {
	query := `SELECT count(*) FROM audios WHERE user_id = ?`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	row := r.db.QueryRow(query, userID)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to scan count: %w", err)
	}

	return count, nil
}

// ListAudiosByUserID returns all audios associated with a given user ID
func (r *mediaRepo) ListAudiosByUserIDAdvance(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity) ([]entity.Audio, error) {
	query := `
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at, title
		FROM audios
		WHERE user_id = ? AND title ILIKE ? and is_deleted = FALSE and status in (?)
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?;
		`
	searchKey = "%" + searchKey + "%"
	query, arg, err := sqlx.In(query, userID, searchKey, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Failed used sqlx in: %v", err)
	}
	query = r.db.Rebind(query)
	rows, err := r.db.Query(query, arg...)
	if err != nil {
		return nil, fmt.Errorf("failed to query audios by user: %v", err)
	}
	defer rows.Close()

	var audios []entity.Audio
	for rows.Next() {
		var a entity.Audio
		var transcriptionID sql.NullInt64
		var videoID sql.NullInt64
		if err := rows.Scan(
			&a.ID,
			&videoID,
			&a.UserID,
			&transcriptionID,
			&a.Duration,
			&a.Lang,
			&a.Folder,
			&a.FileName,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.Title,
		); err != nil {
			return nil, fmt.Errorf("failed to scan audio: %v", err)
		}

		if videoID.Valid {
			a.VideoID = uint64(videoID.Int64)
		} else {
			a.VideoID = 0
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
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at, title
		FROM audios
		WHERE video_id = ? AND id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query) // chuyển ? -> $1, $2

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
		&a.Title,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve audio by VideoID and AudioID: %w", err)
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
		SELECT id, video_id, user_id, transcription_id, duration, lang, folder, file_name, status, created_at, updated_at, title
		FROM audios
		WHERE video_id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := r.db.Query(query, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audios by video: %w", err)
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
			&a.Title,
		); err != nil {
			return nil, fmt.Errorf("failed to scan audio: %w", err)
		}

		if transcriptionID.Valid {
			a.TranscriptionID = uint64(transcriptionID.Int64)
		} else {
			a.TranscriptionID = 0
		}

		audios = append(audios, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating audio rows: %w", err)
	}

	return audios, nil
}


// DeleteAudioByID deletes an audio record by its ID
func (r *mediaRepo) DeleteAudioByID(audioID uint64) error {
	query := "UPDATE audios SET is_deleted = true WHERE id = ?"
	query = sqlx.Rebind(sqlx.DOLLAR, query) // ? → $1

	_, err := r.db.Exec(query, audioID)
	if err != nil {
		return fmt.Errorf("failed to delete audio %d: %w", audioID, err)
	}
	return nil
}


// UpdateAudio updates the entire Audio record
func (r *mediaRepo) UpdateAudio(audio *entity.Audio) error {
	var setClauses []string
	var args []interface{}

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
	if audio.Title != "" {
		setClauses = append(setClauses, "title = ?")
		args = append(args, audio.Title)
	}

	// Always update updated_at
	now := time.Now()
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, now)

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// WHERE clause with audio ID
	args = append(args, audio.ID)
	query := fmt.Sprintf("UPDATE audios SET %s WHERE id = ?", strings.Join(setClauses, ", "))

	// Convert ? to $1, $2,... for PostgreSQL
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

	query = sqlx.Rebind(sqlx.DOLLAR, query) // ? → $1, $2, $3

	now := time.Now()
	result, err := r.db.Exec(query, status, now, audioID)
	if err != nil {
		return fmt.Errorf("failed to update audio status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no audio found with id %d", audioID)
	}

	return nil
}
