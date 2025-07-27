package request

type SetCostRequest struct {
	Pipeline string `json:"pipeline" binding:"required"`
	Model    string `json:"model" binding:"required"`
	Cost     int    `json:"cost" binding:"required"`
}

type SetFeatureFlagActiveRequest struct {
	FlagKey  string `json:"flag_key" binding:"required"`
	IsActive bool   `json:"is_active"`
}
