package media_repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
	"strings"
	"time"
)

func (r *mediaRepo) CreateVideo(video *entity.Video) (uint64, error) {
	if video.Status == "" {
		video.Status = entity.StatusRaw
	}

	// If IDs are zero, we can pass NULL to the DB
	// Otherwise, pass the actual value.
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
		INSERT INTO videos (original_video_id, audio_id, title, duration, description, file_name, folder, image, status, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := r.db.Exec(
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
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert video: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get inserted video ID: %v", err)
	}

	return uint64(id), nil
}

func (r *mediaRepo) GetVideoByID(videoID uint64) (*entity.Video, error) {
	query := `
		SELECT id, original_video_id, audio_id, title, duration, description, file_name, folder, image, status, user_id, created_at, updated_at
		FROM videos
		WHERE id = ?`
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
		return nil, fmt.Errorf("failed to retrieve video: %v", err)
	}

	// If NULL in DB, set 0
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

func (r *mediaRepo) ListVideosByUserID(userID uint64) ([]entity.Video, error) {
	query := `
		SELECT id, original_video_id, audio_id, title, duration, description, file_name, folder, image, status, user_id, created_at, updated_at
		FROM videos
		WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
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
	query := "DELETE FROM videos WHERE id = ?"
	_, err := r.db.Exec(query, videoID)
	if err != nil {
		return fmt.Errorf("failed to delete video %d: %v", videoID, err)
	}
	return nil
}

func (r *mediaRepo) UpdateVideo(video *entity.Video) error {
	var setClauses []string
	var args []interface{}

	// Handle OriginalVideoID
	if video.OriginalVideoID != 0 {
		setClauses = append(setClauses, "original_video_id = ?")
		args = append(args, video.OriginalVideoID)
	}

	// Handle AudioID
	if video.AudioID != 0 {
		setClauses = append(setClauses, "audio_id = ?")
		args = append(args, video.AudioID)
	} else {
		// Set to NULL if zero
		setClauses = append(setClauses, "audio_id = NULL")
	}

	// Handle Title
	if video.Title != "" {
		setClauses = append(setClauses, "title = ?")
		args = append(args, video.Title)
	}

	// Handle Duration
	if video.Duration != 0 {
		setClauses = append(setClauses, "duration = ?")
		args = append(args, video.Duration)
	}

	// Handle Description
	if video.Description != "" {
		setClauses = append(setClauses, "description = ?")
		args = append(args, video.Description)
	}

	// Handle FileName
	if video.FileName != "" {
		setClauses = append(setClauses, "file_name = ?")
		args = append(args, video.FileName)
	}

	// Handle Folder
	if video.Folder != "" {
		setClauses = append(setClauses, "folder = ?")
		args = append(args, video.Folder)
	}

	// Handle Image
	if video.Image != "" {
		setClauses = append(setClauses, "image = ?")
		args = append(args, video.Image)
	}

	// Handle UserID
	if video.UserID != 0 {
		setClauses = append(setClauses, "user_id = ?")
		args = append(args, video.UserID)
	}

	// Always update updated_at
	now := time.Now()
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, now)

	// Check if there are any fields to update
	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add the video.ID for the WHERE clause
	args = append(args, video.ID)

	// Construct the final SQL query
	query := fmt.Sprintf("UPDATE videos SET %s WHERE id = ?", strings.Join(setClauses, ", "))

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
		return fmt.Errorf("no video found with id %d", video.ID)
	}

	return nil
}

func (r *mediaRepo) UpdateVideoStatus(videoID uint64, status entity.StatusEntity) error {
	query := `
		UPDATE videos
		SET status = ?, updated_at = ?
		WHERE id = ?`
	now := time.Now()
	result, err := r.db.Exec(query, status, now, videoID)
	if err != nil {
		return fmt.Errorf("failed to update video status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no video found with id %d", videoID)
	}

	return nil
}

func (r *mediaRepo) GetVideoStatus(videoID uint64) (entity.StatusEntity, error) {
	var status entity.StatusEntity
	query := `SELECT status FROM videos WHERE id = ?`
	err := r.db.QueryRow(query, videoID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("video with ID %d does not exist", videoID)
		}
		return "", fmt.Errorf("failed to get status for video %d: %v", videoID, err)
	}
	return status, nil
}
