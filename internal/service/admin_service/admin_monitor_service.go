package admin_service

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"time"
)

func (s *adminService) GetMonitorDataType(ctx context.Context, adminID uint64) (entity.MonitorDataType, error) {
	validateAdminRole := s.isAdmin(adminID)
	if !validateAdminRole {
		return entity.MonitorDataType{}, fmt.Errorf("only admin can read the model options")
	}

	return s.adminRepo.GetMonitorDataType(ctx)
}

func (s *adminService) GetMonitorPipeline(ctx context.Context, adminID uint64) (entity.MonitorPipeline, error) {
	// check admin role
	if !s.isAdmin(adminID) {
		return entity.MonitorPipeline{}, fmt.Errorf("only admin can retrieve pipeline metrics")
	}
	return s.adminRepo.GetMonitorPipeline(ctx)
}

func (s *adminService) GetMonitorTraffic(
	ctx context.Context,
	adminID uint64,
	periodType entity.TimePeriodType,
	baseTime time.Time,
) (entity.MonitorTraffics, error) {
	// check admin role
	if !s.isAdmin(adminID) {
		return entity.MonitorTraffics{}, fmt.Errorf("only admin can read traffic usage")
	}
	return s.adminRepo.GetMonitorTraffic(ctx, periodType, baseTime)
}
