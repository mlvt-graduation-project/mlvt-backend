package response

import (
	"mlvt/internal/entity"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// JobResponse represents the immediate response for async processing
type JobResponse struct {
	Message string `json:"message"`
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
}

// CallbackRequest represents the callback payload from EC2
type CallbackRequest struct {
	JobID  string      `json:"job_id"`
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// GlobalCallbackURLRequest represents request to set global callback URL
type GlobalCallbackURLRequest struct {
	CallbackURL string `json:"callback_url" binding:"required"`
}

// GlobalCallbackURLResponse represents response with current global callback URL
type GlobalCallbackURLResponse struct {
	CallbackURL string `json:"callback_url"`
}

// StatusResponse represents the response for GetVideoStatus
type StatusResponse struct {
	Status entity.StatusEntity `json:"status"`
}

// MessageResponse represents a message response
type MessageResponse struct {
	Message string `json:"message"`
}

// MessageCreateResponseWithID represents a message response
type MessageCreateResponseWithID struct {
	Message string      `json:"message"`
	Id      interface{} `json:"id"`
}

// TokenResponse represents the response containing a token
type TokenResponse struct {
	Token  string                `json:"token"`
	UserID uint64                `json:"user_id"`
	Role   entity.UserPermission `json:"role"`
}

// AvatarDownloadURLResponse represents the response containing avatar download URL
type AvatarDownloadURLResponse struct {
	AvatarDownloadURL string `json:"avatar_download_url"`
}

// AvatarUploadURLResponse represents the response containing avatar upload URL
type AvatarUploadURLResponse struct {
	AvatarUploadURL string `json:"avatar_upload_url"`
}

// UserResponse represents a single user response
type UserResponse struct {
	User entity.User `json:"user"`
}

// UsersResponse represents multiple users response
type UsersResponse struct {
	Users []entity.User `json:"users"`
}

// UploadURLResponse represents the response containing an upload URL
type UploadURLResponse struct {
	UploadURL string `json:"upload_url"`
}

// DownloadURLResponse represents the response containing a download URL
type DownloadURLResponse struct {
	DownloadURL string `json:"download_url"`
}

// TranscriptionResponse represents the response containing a transcription and its download URL
type TranscriptionResponse struct {
	Transcription entity.Transcription `json:"transcription"`
	DownloadURL   string               `json:"download_url"`
}

// TranscriptionsResponse represents the response containing a list of transcriptions
type TranscriptionsResponse struct {
	Transcriptions []entity.Transcription `json:"transcriptions"`
}

// AudioResponse represents the response containing an audio and its download URL
type AudioResponse struct {
	Audio       entity.Audio `json:"audio"`
	DownloadURL string       `json:"download_url"`
}

// AudiosResponse represents the response containing a list of audios
type AudiosResponse struct {
	Audios []entity.Audio `json:"audios"`
}

type ListVideosByUserIDResponse struct {
	Video    entity.Video `json:"video"`
	VideoURL string       `json:"video_url"`
	ImageURL string       `json:"image_url"`
}

func (v *ListVideosByUserIDResponse) ToProcessResponse() *entity.Process {
	return &entity.Process{
		ID:              strconv.FormatUint(v.Video.ID, 10),
		UserID:          v.Video.UserID,
		MediaType:       "video",
		OriginalVideoID: v.Video.OriginalVideoID,
		Status:          v.Video.Status,
		CreatedAt:       v.Video.CreatedAt,
		UpdatedAt:       v.Video.UpdatedAt,
		ThumbnailUrl:    v.ImageURL,
		VideoUrl:        v.VideoURL,
		Title:           v.Video.Title,
	}
}

type PingStatusResponse struct {
	Status entity.StatusEntity `json:"status"`
}

// Progress with thumbnail
type ProgressResponse struct {
	ID                        primitive.ObjectID  `json:"id"`
	UserID                    uint64              `json:"user_id"`
	ProgressType              entity.ProgressType `json:"progress_type"`
	OriginalVideoID           uint64              `json:"original_video_id"`
	OriginalTranscriptionID   uint64              `json:"original_transcription_id"`
	TranslatedTranscriptionID uint64              `json:"translated_transcription_id"`
	AudioID                   uint64              `json:"audio_id"`
	ProgressedVideoID         uint64              `json:"progressed_video_id"`
	Status                    entity.StatusEntity `json:"status"`
	CreatedAt                 time.Time           `json:"created_at"`
	UpdatedAt                 time.Time           `json:"updated_at"`
	ThumbnailUrl              string              `json:"thumbnail_url"`
	Title                     string              `json:"title"`
}
type ProcessResponse struct {
	TotalCount  int              `json:"total_count"`
	ProcessList []entity.Process `json:"process_list"`
}

func (p *ProgressResponse) ToProcessResponse() *entity.Process {
	return &entity.Process{
		ID:                        p.ID.Hex(),
		UserID:                    p.UserID,
		ProgressType:              p.ProgressType,
		OriginalVideoID:           p.OriginalVideoID,
		OriginalTranscriptionID:   p.OriginalTranscriptionID,
		TranslatedTranscriptionID: p.TranslatedTranscriptionID,
		AudioID:                   p.AudioID,
		ProgressedVideoID:         p.ProgressedVideoID,
		Status:                    p.Status,
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.CreatedAt,
		ThumbnailUrl:              p.ThumbnailUrl,
		Title:                     p.Title,
	}
}

func (p *PingStatusResponse) ValidateStatus() {
	// ToDo: validate string StatusEntity
}

// DailyTokenClaimsResponse represents a list of token claim logs
type DailyTokenClaimsResponse struct {
	Claims []entity.TokenClaim `json:"claims"`
}

// PremiumUsersResponse represents a list of active premium users
type PremiumUsersResponse struct {
	Users []entity.PremiumUser `json:"users"`
}
