package voucher_handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"mlvt/internal/entity"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/voucher_service"
)

type VoucherController struct {
	voucherSvc voucher_service.VoucherService
}

func NewVoucherController(voucherSvc voucher_service.VoucherService) *VoucherController {
	return &VoucherController{
		voucherSvc: voucherSvc,
	}
}

// CreateVoucher godoc
// @Summary Create a new voucher
// @Description Creates a new voucher with specified data.
// @Tags Voucher
// @Accept json
// @Produce json
// @Param voucher body entity.VoucherCode true "Voucher info"
// @Success 200 {object} gin.H
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /voucher [post]
func (vc *VoucherController) CreateVoucher(c *gin.Context) {
	var req entity.VoucherCode
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid request"})
		return
	}

	newID, err := vc.voucherSvc.CreateVoucher(context.Background(), req)
	if err != nil {
		log.Errorf("Failed to create new voucher: %v", err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to create voucher"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "voucher created successfully",
		"id":      newID.Hex(),
	})
}

// UseVoucher godoc
// @Summary Use a voucher by code
// @Description Increments the used count of a voucher code if valid and not expired.
// @Tags Voucher
// @Accept  json
// @Produce  json
// @Param   code path     string true "Voucher Code"
// @Success 200 {object} entity.VoucherCode
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /voucher/use/{code} [post]
func (vc *VoucherController) UseVoucher(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid voucher code"})
		return
	}

	voucher, err := vc.voucherSvc.UseVoucher(context.Background(), code)
	if err != nil {
		log.Errorf("Failed to use voucher: %v", err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, voucher)
}

// UpdateVoucher godoc
// @Summary Updates voucher info
// @Description Updates specific voucher fields. For instance, adjust usage data.
// @Tags Voucher
// @Accept json
// @Produce json
// @Param id path string true "Voucher ID"
// @Param fields body map[string]interface{} true "Fields to update"
// @Success 200 {object} gin.H
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /voucher/{id} [patch]
func (vc *VoucherController) UpdateVoucher(c *gin.Context) {
	idHex := c.Param("id")
	if idHex == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid voucher ID"})
		return
	}

	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid voucher ID format"})
		return
	}

	var fields map[string]interface{}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid request body"})
		return
	}

	if err := vc.voucherSvc.UpdateVoucher(context.Background(), oid, fields); err != nil {
		log.Errorf("Failed to update voucher: %v", err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to update voucher"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "voucher updated successfully"})
}

// GetAllVouchers godoc
// @Summary Retrieves all vouchers
// @Description Returns a list of all vouchers without pagination.
// @Tags Voucher
// @Accept json
// @Produce json
// @Success 200 {array} entity.VoucherCode
// @Failure 500 {object} response.ErrorResponse
// @Router /voucher [get]
func (vc *VoucherController) GetAllVouchers(c *gin.Context) {
	vouchers, err := vc.voucherSvc.GetAllVouchers(context.Background())
	if err != nil {
		log.Errorf("Failed to retrieve vouchers: %v", err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "failed to retrieve vouchers"})
		return
	}

	c.JSON(http.StatusOK, vouchers)
}

// GetVoucherByID godoc
// @Summary Get voucher by ID
// @Description Retrieves a single voucher based on its unique ID
// @Tags Voucher
// @Accept  json
// @Produce  json
// @Param   id path string true "Voucher ID"
// @Success 200 {object} entity.VoucherCode
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /voucher/{id} [get]
func (vc *VoucherController) GetVoucherByID(c *gin.Context) {
	idHex := c.Param("id")
	if idHex == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid voucher ID"})
		return
	}

	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid voucher ID format"})
		return
	}

	// Call service to retrieve voucher
	voucher, err := vc.voucherSvc.GetVoucherByID(context.Background(), oid)
	if err != nil {
		log.Errorf("Failed to get voucher by ID %s: %v", idHex, err)
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Voucher not found"})
		return
	}

	c.JSON(http.StatusOK, voucher)
}
