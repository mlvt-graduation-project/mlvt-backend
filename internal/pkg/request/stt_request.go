package request

// STTRequest represents the request payload for STT processing.
type STTRequest struct {
	InputFileName  string `json:"input_file_name" binding:"required"`
	InputLink      string `json:"input_link" binding:"required"`
	OutputFileName string `json:"output_file_name" binding:"required"`
	OutputLink     string `json:"output_link" binding:"required"`
}
