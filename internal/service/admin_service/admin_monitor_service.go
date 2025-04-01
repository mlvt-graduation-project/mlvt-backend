package admin_service

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
)

func (s *adminService) GetMonitorDataType(ctx context.Context, adminID uint64) (entity.MonitorDataType, error) {
	validateAdminRole := s.isAdmin(adminID)
	if !validateAdminRole {
		return entity.MonitorDataType{}, fmt.Errorf("only admin can read the model options")
	}

	return s.adminRepo.GetMonitorDataType(ctx)
}
