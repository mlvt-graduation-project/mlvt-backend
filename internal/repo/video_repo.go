package repo

import (
	"database/sql"
	"mlvt/internal/entity"
	"time"
)

type VideoRepository interface {
	CreateVideo(video *entity.Video) error
	GetVideoByID(videoID uint64) (*entity.Video, error)
	ListVideosByUserID(userID uint64) ([]entity.Video, error)
	DeleteVideo(videoID uint64) error
	UpdateVideo(video *entity.Video) error
}

type videoRepo struct {
	db *sql.DB
}

func NewVideoRepo(db *sql.DB) VideoRepository {
	return &videoRepo{db: db}
}

// CreateVideo inserts a new video record into the database
func (r *videoRepo) CreateVideo(video *entity.Video) error {
	query := `
		INSERT INTO videos (title, duration, description, file_name, folder, image, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	_, err := r.db.Exec(query, video.Title, video.Duration, video.Description, video.FileName, video.Folder, video.Image, video.UserID, now, now)
	return err
}

// GetVideoByID retrieves a video record by its ID
func (r *videoRepo) GetVideoByID(videoID uint64) (*entity.Video, error) {
	query := `SELECT id, title, duration, description, file_name, folder, image, user_id, created_at, updated_at
	          FROM videos WHERE id = ?`
	row := r.db.QueryRow(query, videoID)
	video := &entity.Video{}
	err := row.Scan(&video.ID, &video.Title, &video.Duration, &video.Description, &video.FileName, &video.Folder, &video.Image, &video.UserID, &video.CreatedAt, &video.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return video, err
}

// ListVideosByUserID lists all videos uploaded by a specific user
func (r *videoRepo) ListVideosByUserID(userID uint64) ([]entity.Video, error) {
	query := `SELECT id, title, duration, description, file_name, folder, image, user_id, created_at, updated_at
	          FROM videos WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []entity.Video
	for rows.Next() {
		var video entity.Video
		if err := rows.Scan(&video.ID, &video.Title, &video.Duration, &video.Description, &video.FileName, &video.Folder, &video.Image, &video.UserID, &video.CreatedAt, &video.UpdatedAt); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

// DeleteVideo deletes a video record by its ID
func (r *videoRepo) DeleteVideo(videoID uint64) error {
	query := "DELETE FROM videos WHERE id = ?"
	_, err := r.db.Exec(query, videoID)
	return err
}

// UpdateVideo updates an existing video record
func (r *videoRepo) UpdateVideo(video *entity.Video) error {
	query := `
		UPDATE videos
		SET title = ?, duration = ?, description = ?, file_name = ?, folder = ?, image = ?, updated_at = ?
		WHERE id = ?`
	now := time.Now()
	_, err := r.db.Exec(query, video.Title, video.Duration, video.Description, video.FileName, video.Folder, video.Image, now, video.ID)
	return err
}
