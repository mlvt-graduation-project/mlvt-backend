package admin_service

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/infra/zap-logging/log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *adminService) GetServerConfig(ctx context.Context) (*entity.AdminConfig, error) {
	return s.adminRepo.GetAdminConfig(ctx)
}

func (s *adminService) UpdateServerConfig(ctx context.Context, adminID uint64, modelType string, modelName string) error {
	if !s.isAdmin(adminID) {
		return fmt.Errorf("access denied: only admins can change the server configuration")
	}

	if !s.isValidField(modelType) {
		return fmt.Errorf("invalid configuration field: %s is not an accepted field", modelType)
	}

	qo := mongodb.QueryOptions{
		Filters: []mongodb.FilterCondition{
			{
				Key:       "model_type",
				Operation: mongodb.OpEqual,
				Value:     modelType,
			},
			{
				Key:       "model_name",
				Operation: mongodb.OpEqual,
				Value:     modelName,
			},
		},
	}

	modelList, err := s.adminRepo.LoadModelOptions(ctx, qo)
	if err != nil {
		return fmt.Errorf("cannot load model with model type: %s and model name: %s, error: %w", modelType, modelName, err)
	}

	if modelList == nil {
		return fmt.Errorf("cannot find your updated config, please add this to model option first")
	}

	filter := bson.M{"config_key": "mlvt"}

	fieldName := s.convertTypeField(modelType)
	updateData := bson.M{
		fieldName:    modelName,
		"updated_at": time.Now(),
	}

	newTraffic := &entity.Traffic{
		ActionType:     entity.AdminDefaultConfigAction,
		Description:    "update current default config",
		UserID:         adminID,
		UserPermission: entity.AdminRole,
	}
	newTraffic.SetCurrentTimestamp()

	_, err = s.trafficService.CreateTraffic(ctx, *newTraffic)
	if err != nil {
		log.Error("failed to log traffic %w", newTraffic)
	}

	return s.adminRepo.UpdateConfig(ctx, filter, updateData)

}
