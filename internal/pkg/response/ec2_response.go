package response

type EC2Response struct {
	Message string `json:"message"`
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Result  string `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
}
