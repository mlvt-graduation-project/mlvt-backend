package schema

import (
	"time"
)

// Video represents the schema for video data
type Video struct {
	ID        uint64    `json:"id"`
	Duration  int       `json:"duration"`   // Duration of the video in seconds
	Link      string    `json:"link"`       // URL to the video on AWS S3
	UserID    uint64    `json:"user_id"`    // ID of the user who uploaded the video
	CreatedAt time.Time `json:"created_at"` // Timestamp of when the video was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp of the last update to the video
}

// AddVideoRequest is used to add a new video link for a user
type AddVideoRequest struct {
	UserID   uint64 `json:"user_id" validate:"required"`
	Link     string `json:"link" validate:"required,url"`
	Duration int    `json:"duration" validate:"required"`
}

// GetVideosResponse represents the response structure for fetching a list of videos
type GetVideosResponse struct {
	Videos []Video `json:"videos"`
}
