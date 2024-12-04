package request

// TTSRequest represents the request payload for TTS processing.
type TTSRequest struct {
	BaseRequest
	Lang string `json:"lang" binding:"required"`
}
