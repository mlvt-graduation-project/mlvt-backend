package admin_service

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/repo/admin_repo"
	"mlvt/internal/repo/user_repo"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type AdminService interface{}

type adminService struct {
	userRepo  user_repo.UserRepository
	adminRepo admin_repo.AdminRepository
}

func NewAminService(
	userRepo user_repo.UserRepository,
	adminRepo admin_repo.AdminRepository,
) AdminService {
	return &adminService{
		userRepo:  userRepo,
		adminRepo: adminRepo,
	}
}

func (s *adminService) isAdmin(id uint64) bool {
	userID, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return false
	}

	if userID.Role != "admin" {
		return false
	}

	return true
}

func (s *adminService) isValidField(fieldName string) bool {
	allowedFields := []string{"stt", "tts", "ttt", "ls"}

	for _, field := range allowedFields {
		if strings.EqualFold(field, fieldName) {
			return true
		}
	}
	return false
}

func (s *adminService) convertTypeField(fieldType string) string {
	switch fieldType {
	case "stt":
		return "stt_model"
	case "ttt":
		return "ttt_model"
	case "tts":
		return "tts_model"
	case "ls":
		return "ls_model"
	default:
		return ""
	}
}

func (s *adminService) GetServerConfig(ctx context.Context) (*entity.AdminConfig, error) {
	return s.adminRepo.GetAdminConfig(ctx)
}

func (s *adminService) GetModelList(ctx context.Context, adminID uint64, qo mongodb.QueryOptions) ([]entity.ModelOption, error) {
	validateAdminRole := s.isAdmin(adminID)
	if !validateAdminRole {
		return nil, fmt.Errorf("only admin can read the model options")
	}

	return s.adminRepo.LoadModelOptions(ctx, qo)
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

	return s.adminRepo.UpdateConfig(ctx, filter, updateData)

}
