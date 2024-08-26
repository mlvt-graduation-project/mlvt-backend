package repo

import (
	"database/sql"
	"mlvt/internal/entity"
	"time"
)

// VideoRepository defines the interface for video repository methods
type VideoRepository interface {
	CreateVideo(video *entity.Video) error
	GetVideoByID(id uint64) (*entity.Video, error)
	GetVideosByUserID(userID uint64) ([]entity.Video, error)
	UpdateVideo(video *entity.Video) error
	DeleteVideo(id uint64) error
}

// VideoRepo implements VideoRepository for working with video data
type VideoRepo struct {
	DB *sql.DB
}

// NewVideoRepo creates a new VideoRepo
func NewVideoRepo(db *sql.DB) *VideoRepo {
	return &VideoRepo{DB: db}
}

// GetVideosByUserID fetches all videos uploaded by a specific user
func (repo *VideoRepo) GetVideosByUserID(userID uint64) ([]entity.Video, error) {
	query := `SELECT id, duration, link, user_id, created_at, updated_at FROM videos WHERE user_id = ?`
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []entity.Video
	for rows.Next() {
		video := entity.Video{}
		if err := rows.Scan(&video.ID, &video.Duration, &video.Link, &video.UserID, &video.CreatedAt, &video.UpdatedAt); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil { // Check for errors encountered during iteration
		return nil, err
	}

	return videos, nil
}

// CreateVideo inserts a new video into the database
func (repo *VideoRepo) CreateVideo(video *entity.Video) error {
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	query := `INSERT INTO videos (duration, link, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err := repo.DB.Exec(query, video.Duration, video.Link, video.UserID, video.CreatedAt, video.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// GetVideoByID fetches a video by its ID
func (repo *VideoRepo) GetVideoByID(id uint64) (*entity.Video, error) {
	video := &entity.Video{}
	query := `SELECT id, duration, link, user_id, created_at, updated_at FROM videos WHERE id = ?`
	err := repo.DB.QueryRow(query, id).Scan(&video.ID, &video.Duration, &video.Link, &video.UserID, &video.CreatedAt, &video.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No video found
		}
		return nil, err
	}
	return video, nil
}

// UpdateVideo updates an existing video's details
func (repo *VideoRepo) UpdateVideo(video *entity.Video) error {
	video.UpdatedAt = time.Now()
	query := `UPDATE videos SET duration = ?, link = ?, user_id = ?, updated_at = ? WHERE id = ?`
	_, err := repo.DB.Exec(query, video.Duration, video.Link, video.UserID, video.UpdatedAt, video.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteVideo removes a video from the database
func (repo *VideoRepo) DeleteVideo(id uint64) error {
	query := `DELETE FROM videos WHERE id = ?`
	_, err := repo.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
