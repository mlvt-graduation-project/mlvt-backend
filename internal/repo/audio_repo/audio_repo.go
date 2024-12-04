package audio_repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
	"time"
)

type AudioRepository interface {
	CreateAudio(audio *entity.Audio) (uint64, error)
	GetAudioByID(audioID uint64) (*entity.Audio, error)
	GetAudioByIDAndUserID(audioID, userID uint64) (*entity.Audio, error)
	ListAudiosByUserID(userID uint64) ([]entity.Audio, error)
	GetAudioByVideoID(videoID, audioID uint64) (*entity.Audio, error)
	ListAudiosByVideoID(videoID uint64) ([]entity.Audio, error)
	DeleteAudioByID(audioID uint64) error
	UpdateAudio(audio *entity.Audio) error
	UpdateAudioStatus(audioID uint64, status entity.StatusEntity) error
}

type audioRepo struct {
	db *sql.DB
}

func NewAudioRepository(db *sql.DB) AudioRepository {
	return &audioRepo{db: db}
}

// CreateAudio inserts a new audio record into the database
func (r *audioRepo) CreateAudio(audio *entity.Audio) (uint64, error) {
	query := `
		INSERT INTO audios (video_id, user_id, duration, lang, folder, file_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		audio.VideoID, audio.UserID, audio.Duration, audio.Lang, audio.Folder, audio.FileName, now, now)
	if err != nil {
		return 0, err
	}

	// Retrieve the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

// GetAudioByID fetches an audio by its ID and user ID
func (r *audioRepo) GetAudioByID(audioID uint64) (*entity.Audio, error) {
	query := `SELECT id, video_id, user_id, duration, lang, folder, file_name, created_at, updated_at
	          FROM audios WHERE id = ?`

	row := r.db.QueryRow(query, audioID)

	audio := &entity.Audio{}
	err := row.Scan(&audio.ID, &audio.VideoID, &audio.UserID, &audio.Duration, &audio.Lang, &audio.Folder,
		&audio.FileName, &audio.CreatedAt, &audio.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return audio, err
}

// GetAudioByIDAndUserID retrieves a single audio by its ID and User ID (owner)
func (r *audioRepo) GetAudioByIDAndUserID(audioID, userID uint64) (*entity.Audio, error) {
	query := `
		SELECT id, video_id, user_id, duration, lang, folder, file_name, created_at, updated_at
		FROM audios WHERE id = ? AND user_id = ?`

	row := r.db.QueryRow(query, audioID, userID)

	audio := &entity.Audio{}
	err := row.Scan(&audio.ID, &audio.VideoID, &audio.UserID, &audio.Duration, &audio.Lang, &audio.Folder, &audio.FileName, &audio.CreatedAt, &audio.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // No record found
	}
	return audio, err
}

// ListAudiosByUserID returns all audios associated with a given user ID
func (r *audioRepo) ListAudiosByUserID(userID uint64) ([]entity.Audio, error) {
	query := `SELECT id, video_id, user_id, duration, lang, folder, file_name, created_at, updated_at
	          FROM audios WHERE user_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audios []entity.Audio
	for rows.Next() {
		var audio entity.Audio
		if err := rows.Scan(&audio.ID, &audio.VideoID, &audio.UserID, &audio.Duration, &audio.Lang, &audio.Folder,
			&audio.FileName, &audio.CreatedAt, &audio.UpdatedAt); err != nil {
			return nil, err
		}
		audios = append(audios, audio)
	}

	return audios, nil
}

// GetAudioByVideoID retrieves a specific audio by its video ID and audio ID
func (r *audioRepo) GetAudioByVideoID(videoID, audioID uint64) (*entity.Audio, error) {
	query := `SELECT id, video_id, user_id, duration, lang, folder, file_name, created_at, updated_at
	          FROM audios WHERE video_id = ? AND id = ?`

	row := r.db.QueryRow(query, videoID, audioID)

	audio := &entity.Audio{}
	err := row.Scan(&audio.ID, &audio.VideoID, &audio.UserID, &audio.Duration, &audio.Lang, &audio.Folder,
		&audio.FileName, &audio.CreatedAt, &audio.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return audio, err
}

// ListAudiosByVideoID returns all audios associated with a given video ID
func (r *audioRepo) ListAudiosByVideoID(videoID uint64) ([]entity.Audio, error) {
	query := `SELECT id, video_id, user_id, duration, lang, folder, file_name, created_at, updated_at
	          FROM audios WHERE video_id = ?`

	rows, err := r.db.Query(query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audios []entity.Audio
	for rows.Next() {
		var audio entity.Audio
		if err := rows.Scan(&audio.ID, &audio.VideoID, &audio.UserID, &audio.Duration, &audio.Lang, &audio.Folder,
			&audio.FileName, &audio.CreatedAt, &audio.UpdatedAt); err != nil {
			return nil, err
		}
		audios = append(audios, audio)
	}

	return audios, nil
}

// DeleteAudioByID deletes an audio record by its ID
func (r *audioRepo) DeleteAudioByID(audioID uint64) error {
	query := "DELETE FROM audios WHERE id = ?"
	_, err := r.db.Exec(query, audioID)
	return err
}

// UpdateAudio updates the entire Audio record
func (r *audioRepo) UpdateAudio(audio *entity.Audio) error {
	query := `
        UPDATE audios
        SET video_id = ?, user_id = ?, duration = ?, lang = ?, folder = ?, file_name = ?, status = ?, updated_at = ?
        WHERE id = ?`
	now := time.Now()
	result, err := r.db.Exec(query,
		audio.VideoID, audio.UserID, audio.Duration, audio.Lang, audio.Folder,
		audio.FileName, audio.Status, now, audio.ID)
	if err != nil {
		return fmt.Errorf("failed to execute update: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no audio found with id %d", audio.ID)
	}
	return nil
}

// UpdateAudioStatus updates only the status of an Audio record
func (r *audioRepo) UpdateAudioStatus(audioID uint64, status entity.StatusEntity) error {
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
