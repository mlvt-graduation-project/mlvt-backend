package request

type BaseRequest struct {
	InputFileName  string `json:"input_file_name" binding:"required"`
	InputLink      string `json:"input_link" binding:"required"`
	OutputFileName string `json:"output_file_name" binding:"required"`
	OutputLink     string `json:"output_link" binding:"required"`
	Model          string `json:"model"`
}

type BaseLang struct {
	SourceLang string `json:"source_language" binding:"required"`
	TargetLang string `json:"target_language" binding:"required"`
}
