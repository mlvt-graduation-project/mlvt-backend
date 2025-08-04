package feature_flag_repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/zap-logging/log"

	"github.com/jmoiron/sqlx"
)

type FeatureFlagRepository interface {
	GetPipelineCost(pipeline string, model string) (int, bool, error)
	SetPipelineCost(pipeline string, model string, cost int) (int, error)
	GetAllFeatureFlag() ([]entity.FeatureFlag, error)
	SetFlagActive(flagKey string, isActive bool) error
	GetFeatureFlagByFlagKey(flagKey string) (*entity.FeatureFlag, error)
}

type featureFlagRepo struct {
	db *sqlx.DB
}

func NewFeatureFlagRepo(db *sqlx.DB) FeatureFlagRepository {
	return &featureFlagRepo{db}
}

func (r *featureFlagRepo) GetPipelineCost(pipeline string, model string) (int, bool, error) {
	var raw entity.FeatureFlag

	query := `SELECT config_details, is_active FROM feature_flags WHERE flag_key = 'model charge' LIMIT 1`
	err := r.db.Get(&raw, query)
	if err != nil {
		log.Errorf("DB error: failed to SELECT config_details for model_charge: %v", err)
		return 0, false, err
	}

	var config map[string]map[string]int
	if err := json.Unmarshal(raw.ConfigDetails, &config); err != nil {
		log.Errorf("JSON unmarshal error: config_details is not valid structure: %v", err)
		return 0, false, err
	}

	if modelCost, ok := config[pipeline]; ok {
		if cost, ok := modelCost[model]; ok {
			return cost, raw.IsActive, nil
		}
	}

	return 0, raw.IsActive, fmt.Errorf("cost not found for pipeline: %s, model: %s", pipeline, model)
}

func (r *featureFlagRepo) SetPipelineCost(pipeline string, model string, cost int) (int, error) {
	var raw entity.FeatureFlag

	tx, err := r.db.Beginx()
	if err != nil {
		log.Errorf("DB error: failed to begin transaction: %v", err)
		return 0, err
	}
	defer tx.Rollback()

	query := `SELECT config_details FROM feature_flags WHERE flag_key = 'model charge' FOR UPDATE`
	err = tx.Get(&raw, query)
	if err != nil {
		log.Errorf("DB error: failed to SELECT FOR UPDATE config_details: %v", err)
		return 0, err
	}

	var config map[string]map[string]int
	if err := json.Unmarshal(raw.ConfigDetails, &config); err != nil {
		log.Errorf("JSON unmarshal error: invalid config_details in SetPipelineCost: %v", err)
		return 0, err
	}

	if _, ok := config[pipeline]; !ok {
		config[pipeline] = make(map[string]int)
	}
	config[pipeline][model] = cost

	newConfig, err := json.Marshal(config)
	if err != nil {
		log.Errorf("JSON marshal error: cannot encode updated config: %v", err)
		return 0, err
	}

	updateQuery := `UPDATE feature_flags SET config_details = $1, updated_at = NOW() WHERE flag_key = 'model charge'`
	_, err = tx.Exec(updateQuery, newConfig)
	if err != nil {
		log.Errorf("DB error: failed to UPDATE config_details: %v", err)
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("DB error: failed to commit transaction: %v", err)
		return 0, err
	}

	return cost, nil
}

func (r *featureFlagRepo) GetAllFeatureFlag() ([]entity.FeatureFlag, error) {
	var flags []entity.FeatureFlag
	err := r.db.Select(&flags, `SELECT * FROM feature_flags`)
	if err != nil {
		log.Errorf("DB error: failed to SELECT all feature flags: %v", err)
	}
	return flags, err
}

func (r *featureFlagRepo) GetFeatureFlagByFlagKey(flagKey string) (*entity.FeatureFlag, error) {
	var flag entity.FeatureFlag
	query := `SELECT * FROM feature_flags WHERE flag_key = $1 LIMIT 1`
	err := r.db.Get(&flag, query, flagKey)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Errorf("DB error: failed to get feature flag by flag_key=%s: %v", flagKey, err)
		return nil, err
	}
	return &flag, nil
}

func (r *featureFlagRepo) SetFlagActive(flagKey string, isActive bool) error {
	query := `UPDATE feature_flags SET is_active = $1, updated_at = NOW() WHERE flag_key = $2`
	_, err := r.db.Exec(query, isActive, flagKey)
	if err != nil {
		log.Errorf("DB error: failed to update is_active for flag_key=%s: %v", flagKey, err)
	}
	return err
}
