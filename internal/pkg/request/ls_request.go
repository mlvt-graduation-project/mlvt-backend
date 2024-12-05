package request

type LSRequest struct {
	InputVideoFileName  string `json:"input_video_file_name" binding:"required"`
	InputVideoLink      string `json:"input_video_link" binding:"required"`
	InputAudioFileName  string `json:"input_audio_file_name" binding:"required"`
	InputAudioLink      string `json:"input_audio_link" binding:"required"`
	OutputVideoFileName string `json:"output_video_file_name" binding:"required"`
	OutputVideoLink     string `json:"output_video_link" binding:"required"`
	Model               string `json:"model"`
}
