package feature_flag_service

import (
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/repo/feature_flag_repo"
	"mlvt/internal/utility"
)

type FeatureFlagService interface {
	SetPipelineCost(pipeline string, model string, cost int) (int, error)
	GetAllFeatureFlags() (map[string]entity.FeatureFlagMapValue, error)
	SetFeatureFlagActive(flagKey string, isActive bool) error
	GetPipelineActiveAndCost(pipeline string, model string) (int, bool, error)
}

type featureFlagService struct {
	repo feature_flag_repo.FeatureFlagRepository
}

func NewFeatureFlagService(repo feature_flag_repo.FeatureFlagRepository) FeatureFlagService {
	return &featureFlagService{repo}
}

func (r *featureFlagService) GetPipelineActiveAndCost(pipeline string, model string) (int, bool, error) {
	if model == "" {
		model = "default"
	}
	return r.repo.GetPipelineCost(pipeline, model)
}

func (r *featureFlagService) SetPipelineCost(pipeline string, model string, cost int) (int, error) {
	return r.repo.SetPipelineCost(pipeline, model, cost)
}

func (r *featureFlagService) GetAllFeatureFlags() (map[string]entity.FeatureFlagMapValue, error) {
	featureFlags, err := r.repo.GetAllFeatureFlag()
	if err != nil {
		return nil, err
	}
	return utility.FeatureFlagsToMap(featureFlags), nil
}

func (s *featureFlagService) SetFeatureFlagActive(flagKey string, isActive bool) error {
	check, err := s.repo.GetFeatureFlagByFlagKey(flagKey)
	if err != nil {
		return err
	}
	if check == nil {
		return fmt.Errorf("invalid flag key")
	}
	return s.repo.SetFlagActive(flagKey, isActive)
}
