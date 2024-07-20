package repository

import (
	"database/sql"
	"log"
	"mlvt/internal/models"
)

type VideoRepository interface {
	Save(video models.Video) error
	GetVideosByUserID(userID string) ([]models.Video, error)
}

type PostgresVideoRepository struct {
	db *sql.DB
}

func NewPostgresVideoRepository(db *sql.DB) *PostgresVideoRepository {
	return &PostgresVideoRepository{
		db: db,
	}
}

func (r *PostgresVideoRepository) Save(video models.Video) error {
	//save to db
	query := `INSERT INTO videos (id, file_path, uploaded_at, size, duration, type, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, video.ID, video.FilePath, video.UploadedAt, video.Size, video.Duration, video.Type, video.UserID)
	if err != nil {
		return err
	}

	log.Printf("Video metadata saved and file uploaded - ID: %s, S3 Path: %s", video.ID, video.FilePath)
	return nil
}

func (r *PostgresVideoRepository) GetVideosByUserID(userID string) ([]models.Video, error) {
	query := `SELECT id, duration, file_path, uploaded_at, size, user_id, type FROM videos WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var video models.Video
		if err := rows.Scan(&video.ID, &video.Duration, &video.FilePath, &video.UploadedAt, &video.Size, &video.UserID, &video.Type); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}
