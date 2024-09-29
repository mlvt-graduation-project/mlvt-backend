package entity

import "time"

type Transcription struct {
	ID        uint64    `json:"id"`
	VideoID   uint64    `json:"video_id"`   // ID of the related video
	UserID    uint64    `json:"user_id"`    // ID of the user who created the transcription
	Text      string    `json:"text"`       // The transcription text
	Lang      string    `json:"lang"`       // Language of the transcription (e.g., "en", "es", etc.)
	Folder    string    `json:"folder"`     // S3 folder or path containing the transcription file
	FileName  string    `json:"file_name"`  // The transcription file name in S3
	CreatedAt time.Time `json:"created_at"` // Timestamp of when the transcription was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp of the last update to the transcription
}
