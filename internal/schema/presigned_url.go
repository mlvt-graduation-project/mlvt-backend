package schema

// GeneratePresignedURLRequest is used to request a pre-signed URL for video uploads
type GeneratePresignedURLRequest struct {
	UserID   uint64 `json:"user_id" validate:"required"`   // ID of the user uploading the video
	FileName string `json:"file_name" validate:"required"` // Name of the video file to be uploaded
	FileType string `json:"file_type" validate:"required"` // Type of the video file (e.g., video/mp4)
	Title    string `json:"title" validate:"required"`     // Title of the video
	Duration int    `json:"duration" validate:"required"`  // Expected duration of the video
}
