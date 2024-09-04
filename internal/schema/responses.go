package schema

// GeneratePresignedURLResponse defines the structure for the successful response
type GeneratePresignedURLResponse struct {
	PresignedUrl string `json:"presignedUrl"`
}

// ErrorResponse defines the structure for error responses
type ErrorResponse struct {
	Error string `json:"error"`
}
