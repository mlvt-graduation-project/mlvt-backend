package admin_service

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/repo/admin_repo"
	"mlvt/internal/repo/user_repo"
	"mlvt/internal/service/traffic_service"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminService interface {
	// Config
	GetServerConfig(ctx context.Context) (*entity.AdminConfig, error)
	UpdateServerConfig(ctx context.Context, adminID uint64, modelType string, modelName string) error

	// Model Options
	GetModelList(ctx context.Context, adminID uint64, qo mongodb.QueryOptions) ([]entity.ModelOption, error)
	AddModelOption(ctx context.Context, adminID uint64, modelOption entity.ModelOption) (primitive.ObjectID, error)
	UpdateModelOption(ctx context.Context, adminID uint64, modelOption entity.ModelOption) error

	// Monitor
	GetMonitorDataType(ctx context.Context, adminID uint64) (entity.MonitorDataType, error)
	GetMonitorPipeline(ctx context.Context, adminID uint64) (entity.MonitorPipeline, error)
	GetMonitorTraffic(ctx context.Context, adminID uint64, periodType entity.TimePeriodType, baseTime time.Time) (entity.MonitorTraffics, error)
}

type adminService struct {
	userRepo       user_repo.UserRepository
	adminRepo      admin_repo.AdminRepository
	trafficService traffic_service.TrafficService
}

func NewAminService(
	userRepo user_repo.UserRepository,
	adminRepo admin_repo.AdminRepository,
	trafficService traffic_service.TrafficService,
) AdminService {
	return &adminService{
		userRepo:       userRepo,
		adminRepo:      adminRepo,
		trafficService: trafficService,
	}
}
