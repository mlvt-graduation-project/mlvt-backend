package admin_service

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/infra/zap-logging/log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *adminService) GetModelList(ctx context.Context, adminID uint64, qo mongodb.QueryOptions) ([]entity.ModelOption, error) {
	validateAdminRole := s.isAdmin(adminID)
	if !validateAdminRole {
		return nil, fmt.Errorf("only admin can read the model options")
	}

	newTraffic := &entity.Traffic{
		ActionType:     entity.AdminModelOptionAction,
		Description:    "request to get model list from admin",
		UserID:         adminID,
		UserPermission: entity.AdminRole,
	}
	newTraffic.SetCurrentTimestamp()

	_, err := s.trafficService.CreateTraffic(ctx, *newTraffic)
	if err != nil {
		log.Error("failed to log traffic %w", newTraffic)
	}

	return s.adminRepo.LoadModelOptions(ctx, qo)
}

func (s *adminService) AddModelOption(
	ctx context.Context,
	adminID uint64,
	modelOption entity.ModelOption,
) (
	primitive.ObjectID,
	error,
) {
	if !s.isAdmin(adminID) {
		return primitive.NilObjectID, fmt.Errorf("access denied: only admins can add new model option")
	}

	// Check if model name exists yet
	qo := mongodb.QueryOptions{
		Filters: []mongodb.FilterCondition{
			{
				Key:       "model_name",
				Operation: mongodb.OpEqual,
				Value:     modelOption.ModelName,
			},
			{
				Key:       "model_type",
				Operation: mongodb.OpEqual,
				Value:     modelOption.ModelType,
			},
		},
	}

	result, err := s.adminRepo.LoadModelOptions(ctx, qo)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error when loading model option data")
	}

	// if err != mongo.ErrNoDocuments {
	// 	return primitive.NilObjectID, fmt.Errorf("this option already exists in system")
	// }
	if result != nil {
		return primitive.NilObjectID, fmt.Errorf("this option already exists in system")
	}

	newTraffic := &entity.Traffic{
		ActionType:     entity.AdminModelOptionAction,
		Description:    "add new model to system",
		UserID:         adminID,
		UserPermission: entity.AdminRole,
	}
	newTraffic.SetCurrentTimestamp()

	_, err = s.trafficService.CreateTraffic(ctx, *newTraffic)
	if err != nil {
		log.Error("failed to log traffic %w", newTraffic)
	}

	return s.adminRepo.AddModelOptions(ctx, modelOption)
}

func (s *adminService) UpdateModelOption(
	ctx context.Context,
	adminID uint64,
	modelOption entity.ModelOption,
) error {
	if !s.isAdmin(adminID) {
		return fmt.Errorf("access denied: only admins can update the model options")
	}

	// verify that the model option to update actually exists
	originalModelOption, err := s.adminRepo.GetModelOptionByID(ctx, modelOption.ID)
	if err != nil {
		return fmt.Errorf("failed to find document by id %s: %w", modelOption.ID, err)
	}

	// if the model name or type has changed, ensure it doesn't conflict with an existing option
	if originalModelOption.ModelName != modelOption.ModelName || originalModelOption.ModelType != modelOption.ModelType {
		qo := mongodb.QueryOptions{
			Filters: []mongodb.FilterCondition{
				{
					Key:       "model_name",
					Operation: mongodb.OpEqual,
					Value:     modelOption.ModelName,
				},
				{
					Key:       "model_type",
					Operation: mongodb.OpEqual,
					Value:     modelOption.ModelType,
				},
				{
					// Exclude the current document from the search
					Key:       "_id",
					Operation: mongodb.OpNotEqual,
					Value:     modelOption.ID,
				},
			},
		}

		existing, err := s.adminRepo.LoadModelOptions(ctx, qo)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("error checking for existing model options %w", err)
		}
		if existing != nil && len(existing) > 0 {
			return fmt.Errorf("model option with the same name and type already exists")
		}
	}

	updatedFields := bson.M{
		"model_name":   modelOption.ModelName,
		"model_type":   modelOption.ModelType,
		"descriptions": modelOption.Description,
		"updated_at":   time.Now(),
	}

	filter := bson.M{"_id": modelOption.ID}

	if err := s.adminRepo.UpdateModelOption(ctx, filter, updatedFields); err != nil {
		return fmt.Errorf("failed to update model option: %w", err)
	}

	newTraffic := &entity.Traffic{
		ActionType:     entity.AdminModelOptionAction,
		Description:    fmt.Sprintf("update current model id: %d", modelOption.ID),
		UserID:         adminID,
		UserPermission: entity.AdminRole,
	}
	newTraffic.SetCurrentTimestamp()

	_, err = s.trafficService.CreateTraffic(ctx, *newTraffic)
	if err != nil {
		log.Error("failed to log traffic %w", newTraffic)
	}

	return nil
}
