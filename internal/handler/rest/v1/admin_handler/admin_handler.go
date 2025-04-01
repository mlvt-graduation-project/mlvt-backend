package admin_handler

import (
	"mlvt/internal/service/admin_service"
)

type AdminController struct {
	adminService admin_service.AdminService
}

func NewAdminController(adminService admin_service.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}
