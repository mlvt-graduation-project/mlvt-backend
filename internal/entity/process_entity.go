package entity

import "time"

type Process struct {
	ID                        string       `json:"id"`
	UserID                    uint64       `json:"user_id"`
	Title                     string       `json:"title"`
	ProgressType              ProgressType `json:"progress_type"`
	MediaType                 MediaType    `json:"media_type"`
	OriginalVideoID           uint64       `json:"original_video_id"`
	OriginalTranscriptionID   uint64       `json:"original_transcription_id"`
	TranslatedTranscriptionID uint64       `json:"translated_transcription_id"`
	AudioID                   uint64       `json:"audio_id"`
	ProgressedVideoID         uint64       `json:"progressed_video_id"`
	Status                    StatusEntity `json:"status"`
	CreatedAt                 time.Time    `json:"created_at"`
	UpdatedAt                 time.Time    `json:"updated_at"`
	ThumbnailUrl              string       `json:"thumbnail_url"`
	VideoUrl                  string       `json:"video_url"`
	Language                  string       `json:"language"`
}
