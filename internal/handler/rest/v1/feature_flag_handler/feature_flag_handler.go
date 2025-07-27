package feature_flag_handler

import (
	"mlvt/internal/pkg/request"
	"mlvt/internal/pkg/response"
	feature_flag_service "mlvt/internal/service/feature_flag"
	"mlvt/internal/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FeatureFlagHandler struct {
	service feature_flag_service.FeatureFlagService
}

func NewFeatureFlagHandler(repo feature_flag_service.FeatureFlagService) *FeatureFlagHandler {
	return &FeatureFlagHandler{repo}
}

var (
	allowedPipelineSet    = utility.GetAllowedPipelines()
	allowedModelMapByPipe = utility.GetAllowedModelsByPipeline()
)

// @Summary Get cost of a specific model in a pipeline
// @Tags feature-flag
// @Produce json
// @Param pipeline query string true "Pipeline name"
// @Param model query string true "Model name"
// @Success 200 {object} object{pipeline=string,model=string,cost=int,model_charge_active=bool}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /feature-flag/pipeline-model-cost [get]
func (r *FeatureFlagHandler) GetPipelineCost(c *gin.Context) {
	pipeline := c.Query("pipeline")
	model := c.Query("model")
	if model == "" {
		model = "default"
	}
	if !utility.IsValidPipeline(pipeline, allowedPipelineSet) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid pipeline name"})
		return
	}
	if !utility.IsValidModel(pipeline, model, allowedModelMapByPipe) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid model name"})
		return
	}

	cost, isActive, err := r.service.GetPipelineActiveAndCost(pipeline, model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to get cost"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pipeline":            pipeline,
		"model":               model,
		"cost":                cost,
		"model_charge_active": isActive,
	})
}

// @Summary Update cost of a specific model in a pipeline
// @Tags feature-flag
// @Accept json
// @Produce json
// @Param body body request.SetCostRequest true "Pipeline and model to update cost"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /feature-flag/pipeline-model [patch]
func (r *FeatureFlagHandler) SetPipelineCost(c *gin.Context) {
	var req request.SetCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request body"})
		return
	}

	if !utility.IsValidPipeline(req.Pipeline, allowedPipelineSet) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid pipeline name"})
		return
	}
	if !utility.IsValidModel(req.Pipeline, req.Model, allowedModelMapByPipe) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid model name"})
		return
	}

	_, err := r.service.SetPipelineCost(req.Pipeline, req.Model, req.Cost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to set cost"})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Cost updated successfully"})
}

// @Summary Get all feature flags with active status and config details
// @Tags feature-flag
// @Produce json
// @Success 200 {object} map[string]utility.FeatureFlagMapValue
// @Failure 500 {object} response.ErrorResponse
// @Router /feature-flag/get-all [get]
func (r *FeatureFlagHandler) GetAllFeatureFlags(c *gin.Context) {
	flags, err := r.service.GetAllFeatureFlags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to fetch feature flags"})
		return
	}
	c.JSON(http.StatusOK, flags)
}

// @Summary Set active status of a feature flag
// @Tags feature-flag
// @Accept json
// @Produce json
// @Param body body request.SetFeatureFlagActiveRequest true "Feature flag status update"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /feature-flag/config [put]
func (h *FeatureFlagHandler) SetFeatureFlagActive(c *gin.Context) {
	var req request.SetFeatureFlagActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request body"})
		return
	}

	err := h.service.SetFeatureFlagActive(req.FlagKey, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to update feature flag. Detail error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Feature flag updated successfully"})
}
