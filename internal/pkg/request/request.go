package request

type UpdateServerConfigRequest struct {
	ModelType string `json:"model_type" binding:"required"`
	ModelName string `json:"model_name" binding:"required"`
}
