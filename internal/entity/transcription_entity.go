package entity

import (
	"time"
)

type Transcription struct {
	ID                      uint64       `json:"id" db:"id"`
	VideoID                 uint64       `json:"video_id" db:"video_id"`
	UserID                  uint64       `json:"user_id" db:"user_id"`
	OriginalTranscriptionID uint64       `json:"original_transcription_id" db:"original_transcription_id"`
	Text                    string       `json:"text" db:"text"`
	Lang                    string       `json:"lang" db:"lang"`
	Folder                  string       `json:"folder" db:"folder"`
	FileName                string       `json:"file_name" db:"file_name"`
	Status                  StatusEntity `json:"status" db:"status"`
	CreatedAt               time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time    `json:"updated_at" db:"updated_at"`
	Title                   string       `json:"title" db:"title"`
	IsDeleted               bool         `json:"is_deleted" db:"is_deleted"`
}
