package media_repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func (r *mediaRepo) CreateVideo(video *entity.Video) (uint64, error) {
	if video.Status == "" {
		video.Status = entity.StatusRaw
	}

	var originalVideoID interface{}
	if video.OriginalVideoID == 0 {
		originalVideoID = nil
	} else {
		originalVideoID = video.OriginalVideoID
	}

	var audioID interface{}
	if video.AudioID == 0 {
		audioID = nil
	} else {
		audioID = video.AudioID
	}

	query := `
		INSERT INTO videos (
			original_video_id, audio_id, title, duration, description,
			file_name, folder, image, status, user_id, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	now := time.Now()
	var insertedID uint64
	err := r.db.QueryRow(
		query,
		originalVideoID,
		audioID,
		video.Title,
		video.Duration,
		video.Description,
		video.FileName,
		video.Folder,
		video.Image,
		video.Status,
		video.UserID,
		now,
		now,
	).Scan(&insertedID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert video: %w", err)
	}

	return insertedID, nil
}

func (r *mediaRepo) GetVideoByID(videoID uint64) (*entity.Video, error) {
	query := `
		SELECT id, original_video_id, audio_id, title, duration, description, file_name, folder, image, status, user_id, created_at, updated_at
		FROM videos
		WHERE id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	row := r.db.QueryRow(query, videoID)

	var (
		originalVideoID sql.NullInt64
		audioID         sql.NullInt64
		video           entity.Video
	)

	err := row.Scan(
		&video.ID,
		&originalVideoID,
		&audioID,
		&video.Title,
		&video.Duration,
		&video.Description,
		&video.FileName,
		&video.Folder,
		&video.Image,
		&video.Status,
		&video.UserID,
		&video.CreatedAt,
		&video.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve video: %w", err)
	}

	if originalVideoID.Valid {
		video.OriginalVideoID = uint64(originalVideoID.Int64)
	} else {
		video.OriginalVideoID = 0
	}

	if audioID.Valid {
		video.AudioID = uint64(audioID.Int64)
	} else {
		video.AudioID = 0
	}

	return &video, nil
}

func (r *mediaRepo) GetCountVideosByUserId(userID uint64) (int, error) {
	query := `SELECT count(*) FROM videos WHERE user_id = ?`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	row := r.db.QueryRow(query, userID)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to scan count: %w", err)
	}

	return count, nil
}

func (r *mediaRepo) ListVideosByUserID(userID uint64) ([]entity.Video, error) {
	query := `
		SELECT id, original_video_id, audio_id, title, duration, description, file_name, folder, image, status, user_id, created_at, updated_at
		FROM videos
		WHERE user_id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos by user: %w", err)
	}
	defer rows.Close()

	var videos []entity.Video
	for rows.Next() {
		var v entity.Video
		var originalVideoID sql.NullInt64
		var audioID sql.NullInt64

		if err := rows.Scan(
			&v.ID,
			&originalVideoID,
			&audioID,
			&v.Title,
			&v.Duration,
			&v.Description,
			&v.FileName,
			&v.Folder,
			&v.Image,
			&v.Status,
			&v.UserID,
			&v.CreatedAt,
			&v.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}

		if originalVideoID.Valid {
			v.OriginalVideoID = uint64(originalVideoID.Int64)
		} else {
			v.OriginalVideoID = 0
		}
		if audioID.Valid {
			v.AudioID = uint64(audioID.Int64)
		} else {
			v.AudioID = 0
		}

		videos = append(videos, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over video rows: %w", err)
	}

	return videos, nil
}

func (r *mediaRepo) ListVideosByUserIDAdvance(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity) ([]entity.Video, error) {
	query := `
		SELECT id, original_video_id, audio_id, title, duration, description, file_name, folder, image, status, user_id, created_at, updated_at
		FROM videos
		WHERE user_id = ? AND is_deleted = FALSE AND title ILIKE ? and status IN (?)
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?
	`
	searchKey = "%" + searchKey + "%"
	query, arg, err := sqlx.In(query, userID, searchKey, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to build in clause: %v", err)
	}
	query = r.db.Rebind(query)
	rows, err := r.db.Query(query, arg...)
	if err != nil {
		return nil, fmt.Errorf("failed to query videos by user: %v", err)
	}
	defer rows.Close()

	var videos []entity.Video
	for rows.Next() {
		var v entity.Video
		var originalVideoID sql.NullInt64
		var audioID sql.NullInt64

		if err := rows.Scan(
			&v.ID,
			&originalVideoID,
			&audioID,
			&v.Title,
			&v.Duration,
			&v.Description,
			&v.FileName,
			&v.Folder,
			&v.Image,
			&v.Status,
			&v.UserID,
			&v.CreatedAt,
			&v.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan video: %v", err)
		}

		// If NULL in DB, set 0
		if originalVideoID.Valid {
			v.OriginalVideoID = uint64(originalVideoID.Int64)
		} else {
			v.OriginalVideoID = 0
		}
		if audioID.Valid {
			v.AudioID = uint64(audioID.Int64)
		} else {
			v.AudioID = 0
		}

		videos = append(videos, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over video rows: %v", err)
	}

	return videos, nil
}

func (r *mediaRepo) DeleteVideo(videoID uint64) error {
	query := "UPDATE videos SET is_deleted = true WHERE id = ?"
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	_, err := r.db.Exec(query, videoID)
	if err != nil {
		return fmt.Errorf("failed to delete video %d: %w", videoID, err)
	}
	return nil
}

func (r *mediaRepo) UpdateVideo(video *entity.Video) error {
	var setClauses []string
	var args []interface{}

	if video.OriginalVideoID != 0 {
		setClauses = append(setClauses, "original_video_id = ?")
		args = append(args, video.OriginalVideoID)
	}

	if video.AudioID != 0 {
		setClauses = append(setClauses, "audio_id = ?")
		args = append(args, video.AudioID)
	} else {
		setClauses = append(setClauses, "audio_id = NULL")
	}

	if video.Title != "" {
		setClauses = append(setClauses, "title = ?")
		args = append(args, video.Title)
	}

	if video.Duration != 0 {
		setClauses = append(setClauses, "duration = ?")
		args = append(args, video.Duration)
	}

	if video.Description != "" {
		setClauses = append(setClauses, "description = ?")
		args = append(args, video.Description)
	}

	if video.FileName != "" {
		setClauses = append(setClauses, "file_name = ?")
		args = append(args, video.FileName)
	}

	if video.Folder != "" {
		setClauses = append(setClauses, "folder = ?")
		args = append(args, video.Folder)
	}

	if video.Image != "" {
		setClauses = append(setClauses, "image = ?")
		args = append(args, video.Image)
	}

	if video.UserID != 0 {
		setClauses = append(setClauses, "user_id = ?")
		args = append(args, video.UserID)
	}

	// Always update updated_at
	now := time.Now()
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, now)

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add WHERE clause with id
	args = append(args, video.ID)
	query := fmt.Sprintf("UPDATE videos SET %s WHERE id = ?", strings.Join(setClauses, ", "))

	// Rebind ? -> $1, $2... for PostgreSQL
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
		return fmt.Errorf("no video found with id %d", video.ID)
	}

	return nil
}

func (r *mediaRepo) UpdateVideoStatus(videoID uint64, status entity.StatusEntity) error {
	query := `
		UPDATE videos
		SET status = ?, updated_at = ?
		WHERE id = ?`

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	now := time.Now()
	result, err := r.db.Exec(query, status, now, videoID)
	if err != nil {
		return fmt.Errorf("failed to update video status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no video found with id %d", videoID)
	}

	return nil
}

func (r *mediaRepo) GetVideoStatus(videoID uint64) (entity.StatusEntity, error) {
	var status entity.StatusEntity
	query := `SELECT status FROM videos WHERE id = ?`
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	err := r.db.QueryRow(query, videoID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("video with ID %d does not exist", videoID)
		}
		return "", fmt.Errorf("failed to get status for video %d: %w", videoID, err)
	}
	return status, nil
}