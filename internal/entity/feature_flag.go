package entity

import (
	"encoding/json"
	"time"
)

type STTModel struct {
	Default int `json:"default"`
}

type TTTModel struct {
	Default int `json:"default"`
}

type TTSModel struct {
	Default int `json:"default"`
}

type LipSyncModel struct {
	Default int `json:"default"`
}

type FullPipelineModel struct {
	Default int `json:"default"`
}

// use for feature flag model_charge in feature flag table
type ModelCostConfig struct {
	STTConfig          STTModel          `json:"STT"`
	TTTConfig          TTTModel          `json:"TTT"`
	TTSConfig          TTSModel          `json:"TTS"`
	LipSyncConfig      LipSyncModel      `json:"LS"`
	FullPipelineConfig FullPipelineModel `json:"FP"`
}

type FeatureFlag struct {
	ID            int             `db:"id" json:"id"`
	FlagKey       string          `db:"flag_key" json:"flag_key"`
	ConfigDetails json.RawMessage `db:"config_details" json:"config_details"`
	Description   string          `db:"description" json:"description"`
	IsActive      bool            `db:"is_active" json:"is_active"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at" json:"updated_at"`
}

type FeatureFlagMapValue struct {
	IsActive      bool            `json:"is_active"`
	ConfigDetails json.RawMessage `json:"config_details"`
}
