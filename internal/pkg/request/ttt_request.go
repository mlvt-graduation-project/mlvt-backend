package request

// TTTRequest represents the request payload for TTT processing.
type TTTRequest struct {
	InputFileName  string `json:"input_file_name" binding:"required"`
	InputLink      string `json:"input_link" binding:"required"`
	OutputFileName string `json:"output_file_name" binding:"required"`
	OutputLink     string `json:"output_link" binding:"required"`
	SourceLang     string `json:"source_language" binding:"required"`
	TargetLang     string `json:"target_language" binding:"required"`
}
