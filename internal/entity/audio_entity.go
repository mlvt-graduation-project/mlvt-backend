package entity

import "time"

type Audio struct {
	ID        uint64    `json:"id"`
	VideoID   uint64    `json:"video_id"`   // ID of the related video
	UserID    uint64    `json:"user_id"`    // ID of the user who uploaded the audio
	Duration  int       `json:"duration"`   // Duration of the audio in seconds
	Lang      string    `json:"lang"`       // Language of the audio (e.g., "en", "es", etc.)
	Folder    string    `json:"folder"`     // S3 folder or path containing the audio file
	FileName  string    `json:"file_name"`  // The audio file name in S3
	CreatedAt time.Time `json:"created_at"` // Timestamp of when the audio was uploaded
	UpdatedAt time.Time `json:"updated_at"` // Timestamp of the last update to the audio
}
