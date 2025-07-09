package entity

import "time"

type Video struct {
	ID              uint64       `json:"id" db:"id"`
	OriginalVideoID uint64       `json:"original_video_id" db:"original_video_id"`
	AudioID         uint64       `json:"audio_id" db:"audio_id"`
	Title           string       `json:"title" db:"title"`
	Duration        int          `json:"duration" db:"duration"`
	Description     string       `json:"description" db:"description"`
	FileName        string       `json:"file_name" db:"file_name"`
	Folder          string       `json:"folder" db:"folder"`
	Image           string       `json:"image" db:"image"`
	Status          StatusEntity `json:"status" db:"status"`
	UserID          uint64       `json:"user_id" db:"user_id"`
	CreatedAt       time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at" db:"updated_at"`
	IsDeleted       bool         `json:"is_deleted" db:"is_deleted"`
}
