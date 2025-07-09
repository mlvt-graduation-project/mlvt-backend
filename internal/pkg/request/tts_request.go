package request

// TTSRequest represents the request payload for TTS processing.
type TTSRequest struct {
	BaseRequest
	InputAudioFileName string `json:"input_audio_file_name" binding:"required"`
	InputAudioLink     string `json:"input_audio_link" binding:"required"`
	Lang               string `json:"lang" binding:"required"`
}
