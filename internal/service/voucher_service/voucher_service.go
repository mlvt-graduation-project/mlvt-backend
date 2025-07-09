package voucher_service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mlvt/internal/entity"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/repo/voucher_repo"
	"mlvt/internal/service/traffic_service"
	"mlvt/internal/service/wallet_service"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VoucherService interface {
	CreateVoucher(ctx context.Context, vc entity.VoucherCode) (primitive.ObjectID, error)
	UseVoucher(ctx context.Context, code string, userID uint64) (*entity.VoucherCode, error)
	UpdateVoucher(ctx context.Context, voucher entity.VoucherCode) error
	GetAllVouchers(ctx context.Context, req entity.GetAllVoucherRequest) ([]entity.VoucherCode, int64, error)
	GetVoucherByID(ctx context.Context, id primitive.ObjectID) (*entity.VoucherCode, error)
}

type voucherService struct {
	repo           voucher_repo.VoucherRepository
	trafficService traffic_service.TrafficService
	walletService  wallet_service.WalletService
}

func NewVoucherService(
	repo voucher_repo.VoucherRepository,
	trafficService traffic_service.TrafficService,
	walletService wallet_service.WalletService,
) VoucherService {
	return &voucherService{
		repo:           repo,
		trafficService: trafficService,
		walletService:  walletService,
	}
}

// CreateVoucher inserts a new voucher into the database.
func (s *voucherService) CreateVoucher(ctx context.Context, vc entity.VoucherCode) (primitive.ObjectID, error) {
	// Attempt to find an existing voucher with the same code
	existingVoucher, err := s.repo.FindByCode(ctx, vc.Code)

	if err == nil && existingVoucher != nil {
		// We found an existing code. Check if it's expired.
		if time.Now().Before(existingVoucher.ExpiredTime) {
			// This means the code is still valid (not expired).
			return primitive.NilObjectID, fmt.Errorf("voucher code already exists and is not expired")
		}
		// If we reach here, the code exists but is expired, so we can reuse the code.
	}
	// If err != nil, it's possible the code wasn't found (or a database error occurred).
	// We'll proceed to create a new voucher if the code isn't actively valid.

	// Now create the new voucher
	vc.CreatedAt = time.Now()
	vc.UpdatedAt = vc.CreatedAt

	insertedID, insertErr := s.repo.Insert(ctx, vc)
	if insertErr != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert voucher: %w", insertErr)
	}

	// log traffic
	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.AdminVoucherAction,
		Description: fmt.Sprintf("admin create new voucher, id: %s", insertedID.Hex()),
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: deposit to user account, err: %s", err)
	}

	return insertedID, nil
}

// UseVoucher checks if a voucher is valid, updates usage if it is, and returns the updated voucher.
func (s *voucherService) UseVoucher(ctx context.Context, code string, userID uint64) (*entity.VoucherCode, error) {
	voucher, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("voucher not found: %w", err)
	}

	// Check if voucher has expired
	if time.Now().After(voucher.ExpiredTime) {
		return nil, fmt.Errorf("voucher has expired")
	}

	// Check usage limit
	if voucher.UsedCount >= voucher.MaxUsage {
		return nil, fmt.Errorf("voucher usage limit reached")
	}

	// Increase usage count
	voucher.UsedCount++
	voucher.UpdatedAt = time.Now()

	// Update in DB
	filter := bson.M{"_id": voucher.Id}
	updatedFields := bson.M{}
	if voucher.UsedCount != 0 {
		updatedFields["used_count"] = voucher.UsedCount
	}
	updatedFields["updated_at"] = time.Now()
	if err := s.repo.UpdateVoucher(ctx, filter, updatedFields); err != nil {
		return nil, fmt.Errorf("failed to update voucher usage: %w", err)
	}

	err = s.walletService.Deposit(ctx, userID, int64(voucher.Token))
	if err != nil {
		return nil, fmt.Errorf("failed to deposit voucher amount to user wallet: %w", err)
	}

	// log traffic
	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.RedeemVoucherAction,
		Description: fmt.Sprintf("user redeem voucher, id: %s", voucher.Id.Hex()),
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: create new voucher, err: %s", err)
	}

	return voucher, nil
}

// UpdateVoucher updates one or more fields for a given voucher ID.
func (s *voucherService) UpdateVoucher(ctx context.Context, voucher entity.VoucherCode) error {
	existing, err := s.repo.FindByID(ctx, voucher.Id)
	if err != nil {
		return fmt.Errorf("failed to find voucher: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("voucher does not exist")
	}

	updatedFields := bson.M{}
	if voucher.Code != "" {
		updatedFields["code"] = voucher.Code
	}
	if voucher.Token != 0 {
		updatedFields["token"] = voucher.Token
	}
	if voucher.MaxUsage != 0 {
		updatedFields["max_usage"] = voucher.MaxUsage
	}
	if voucher.UsedCount != 0 {
		updatedFields["used_count"] = voucher.UsedCount
	}
	if !voucher.ExpiredTime.IsZero() {
		updatedFields["expired_time"] = voucher.ExpiredTime
	}

	// Always update updated_at
	updatedFields["updated_at"] = time.Now()

	filter := bson.M{"_id": voucher.Id}

	if err := s.repo.UpdateVoucher(ctx, filter, updatedFields); err != nil {
		return fmt.Errorf("failed to update voucher: %w", err)
	}

	// log traffic
	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.AdminVoucherAction,
		Description: fmt.Sprintf("admin update voucher, id: %s", voucher.Id.Hex()),
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: update voucher, err: %s", err)
	}

	return nil
}

// GetAllVouchers returns a list of all vouchers.
func (s *voucherService) GetAllVouchers(ctx context.Context, req entity.GetAllVoucherRequest) ([]entity.VoucherCode, int64, error) {
	sortField := req.SortBy
	if req.SortBy == "" {
		sortField = "EXPIRED_TIME"
	} 
	sortOrder := -1
	if req.Sort == "ASC" {
		sortOrder = 1
	}
	searchField := req.SearchCriteria
	status := strings.ToUpper(req.Status)

	// Log + repo
	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.AdminVoucherAction,
		Description: "admin get all vouchers",
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: get all vouchers, err: %s", err)
	}
	return s.repo.GetAll(ctx, status, strings.ToLower(sortField), sortOrder, strings.ToLower(searchField), req.SearchKey, req.Offset, req.Limit)
}

func (s *voucherService) GetVoucherByID(ctx context.Context, id primitive.ObjectID) (*entity.VoucherCode, error) {
	voucher, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find voucher by ID: %w", err)
	}

	// log traffic
	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.AdminVoucherAction,
		Description: fmt.Sprintf("admin get a voucher, id: %s", id.Hex()),
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: get voucher by id, err: %s", err)
	}
	return voucher, nil
}
