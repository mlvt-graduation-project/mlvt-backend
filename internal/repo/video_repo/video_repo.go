package video_repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
	"time"
)

type VideoRepository interface {
	CreateVideo(video *entity.Video) (uint64, error)
	GetVideoByID(videoID uint64) (*entity.Video, error)
	ListVideosByUserID(userID uint64) ([]entity.Video, error)
	DeleteVideo(videoID uint64) error
	UpdateVideo(video *entity.Video) error
	GetVideoStatus(videoID uint64) (entity.StatusEntity, error)
	UpdateVideoStatus(videoId uint64, status entity.StatusEntity) error
}

type videoRepo struct {
	db *sql.DB
}

func NewVideoRepo(db *sql.DB) VideoRepository {
	return &videoRepo{db: db}
}

// CreateVideo inserts a new video record into the database
func (r *videoRepo) CreateVideo(video *entity.Video) (uint64, error) {
	if video.Status == "" {
		video.Status = entity.StatusRaw
	}
	query := `
		INSERT INTO videos (original_video_id, audio_id, title, duration, description, file_name, folder, image, status, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := r.db.Exec(
		query,
		video.OriginalVideoID,
		video.AudioID,
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
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

// GetVideoByID retrieves a video record by its ID
func (r *videoRepo) GetVideoByID(videoID uint64) (*entity.Video, error) {
	query := `
		SELECT id, original_video_id, audio_id, title, duration, description, file_name, folder, image, status, user_id, created_at, updated_at
		FROM videos
		WHERE id = ?`
	row := r.db.QueryRow(query, videoID)
	video := &entity.Video{}
	err := row.Scan(
		&video.ID,
		&video.OriginalVideoID,
		&video.AudioID,
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
	return video, nil
}

// ListVideosByUserID lists all videos uploaded by a specific user
func (r *videoRepo) ListVideosByUserID(userID uint64) ([]entity.Video, error) {
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
		if err := rows.Scan(
			&v.ID,
			&v.OriginalVideoID,
			&v.AudioID,
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
		videos = append(videos, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating video rows: %v", err)
	}

	return videos, nil
}

// DeleteVideo deletes a video record by its ID
func (r *videoRepo) DeleteVideo(videoID uint64) error {
	query := "DELETE FROM videos WHERE id = ?"
	_, err := r.db.Exec(query, videoID)
	if err != nil {
		return fmt.Errorf("failed to delete video %d: %v", videoID, err)
	}
	return nil
}

// UpdateVideo updates an existing video record
func (r *videoRepo) UpdateVideo(video *entity.Video) error {
	query := `
        UPDATE videos
        SET original_video_id = ?, audio_id = ?, title = ?, duration = ?, description = ?, file_name = ?, folder = ?, image = ?, status = ?, updated_at = ?
        WHERE id = ?`
	now := time.Now()
	result, err := r.db.Exec(query,
		video.OriginalVideoID,
		video.AudioID,
		video.Title,
		video.Duration,
		video.Description,
		video.FileName,
		video.Folder,
		video.Image,
		video.Status,
		now,
		video.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to execute update: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no video found with id %d", video.ID)
	}
	return nil
}

// UpdateVideoStatus updates only the status of a video record
func (r *videoRepo) UpdateVideoStatus(videoID uint64, status entity.StatusEntity) error {
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

func (r *videoRepo) GetVideoStatus(videoID uint64) (entity.StatusEntity, error) {
	var status entity.StatusEntity
	query := `
		SELECT status
		FROM videos
		WHERE id = ?
	`
	err := r.db.QueryRow(query, videoID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("video with ID %d does not exist", videoID)
		}
		return "", fmt.Errorf("failed to get status for video %d: %v", videoID, err)
	}
	return status, nil
}
