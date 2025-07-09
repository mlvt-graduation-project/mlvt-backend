package voucher_handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"mlvt/internal/entity"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/voucher_service"
	"mlvt/internal/utils"
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
// @Success 200 {object} response.MessageCreateResponseWithID
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /voucher/create [post]
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

	c.JSON(http.StatusOK, response.MessageCreateResponseWithID{
		Message: "voucher created successfully",
		Id:      newID.Hex(),
	})
}

type UseVoucherRequest struct {
	UserID uint64 `json:"user_id"`
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

	var req UseVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid request body"})
		return
	}

	voucher, err := vc.voucherSvc.UseVoucher(context.Background(), code, req.UserID)
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
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /voucher/{voucherID} [patch]
func (vc *VoucherController) UpdateVoucher(c *gin.Context) {
	ctx := context.Background()

	voucherIDStr := c.Param("voucherID")
	objectID, err := primitive.ObjectIDFromHex(voucherIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid voucher ID"})
		return
	}

	// Bind request body into a voucher struct.
	var voucher entity.VoucherCode
	if err := c.ShouldBindJSON(&voucher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Ensure we set the correct ID and update timestamp.
	voucher.Id = objectID
	voucher.UpdatedAt = time.Now()

	// Pass along to the service for the actual update.
	if err := vc.voucherSvc.UpdateVoucher(ctx, voucher); err != nil {
		log.Errorf("Failed to update voucher: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to update voucher",
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Message: "voucher updated successfully",
	})
}

// GetAllVouchers godoc
// @Summary Retrieves all vouchers
// @Description Returns a list of all vouchers without pagination.
// @Tags Voucher
// @Accept json
// @Produce json
// @Success 200 {array} entity.VoucherCode
// @Failure 500 {object} response.ErrorResponse
// @Router /voucher/get-all [get]
func (vc *VoucherController) GetAllVouchers(c *gin.Context) {
	var req entity.GetAllVoucherRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid query parameters"})
		return
	}

	// Validate allowed values (can move to helper)
	if req.Status != "" && !utils.IsInListString(req.Status, []string{"ACTIVE", "EXPIRED", "USED"}) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid status"})
		return
	}
	if req.SortBy != "" && !utils.IsInListString(req.SortBy, []string{"ID", "CODE", "TOKEN", "MAX_USAGE", "USED_COUNT", "EXPIRED_TIME", "CREATED_AT", "UPDATED_AT"}) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid sortby"})
		return
	}
	if req.Sort != "" && !utils.IsInListString(req.Sort, []string{"ASC", "DESC"}) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid sort"})
		return
	}
	if req.SearchCriteria != "" && !utils.IsInListString(req.SearchCriteria, []string{"CODE", "TOKEN", "MAX_USAGE", "USED_COUNT", "EXPIRED_TIME"}) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid searchCriteria"})
		return
	}
	if utils.IsInListString(req.SearchCriteria, []string{"TOKEN", "MAX_USAGE", "USED_COUNT"}) {
		if _, err := strconv.Atoi(req.SearchKey); err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "search key must be number"})
			return
		}
	}

	// Gọi service
	vouchers, totalCount, err := vc.voucherSvc.GetAllVouchers(context.Background(), req)
	resp := entity.GetAllVoucherResponse {
		Vouchers: vouchers,
		TotalCount: totalCount,
	}
	if err != nil {
		log.Errorf("Failed to retrieve vouchers: %v", err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
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
// @Router /voucher/{voucherID} [get]
func (vc *VoucherController) GetVoucherByID(c *gin.Context) {
	idHex := c.Param("voucherID")
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
