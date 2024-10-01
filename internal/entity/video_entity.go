package entity

import "time"

type Video struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Duration    int       `json:"duration"` // Duration of the video in seconds
	Description string    `json:"description"`
	FileName    string    `json:"file_name"` //
	Folder      string    `json:"folder"`
	Image       string    `json:"image"`
	UserID      uint64    `json:"user_id"`    // ID of the user who uploaded the video
	CreatedAt   time.Time `json:"created_at"` // Timestamp of when the video was created
	UpdatedAt   time.Time `json:"updated_at"` // Timestamp of the last update to the video
}
