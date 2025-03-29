package voucher_service

import (
	"context"
	"fmt"
	"time"

	"mlvt/internal/entity"
	"mlvt/internal/repo/voucher_repo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VoucherService interface {
	CreateVoucher(ctx context.Context, vc entity.VoucherCode) (primitive.ObjectID, error)
	UseVoucher(ctx context.Context, code string) (*entity.VoucherCode, error)
	UpdateVoucher(ctx context.Context, id primitive.ObjectID, fields map[string]interface{}) error
	GetAllVouchers(ctx context.Context) ([]entity.VoucherCode, error)
}

type voucherService struct {
	repo voucher_repo.VoucherRepository
}

func NewVoucherService(repo voucher_repo.VoucherRepository) VoucherService {
	return &voucherService{
		repo: repo,
	}
}

// CreateVoucher inserts a new voucher into the database.
func (s *voucherService) CreateVoucher(ctx context.Context, vc entity.VoucherCode) (primitive.ObjectID, error) {
	vc.CreatedAt = time.Now()
	vc.UpdatedAt = time.Now()

	return s.repo.Insert(ctx, vc)
}

// UseVoucher checks if a voucher is valid, updates usage if it is, and returns the updated voucher.
func (s *voucherService) UseVoucher(ctx context.Context, code string) (*entity.VoucherCode, error) {
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
	updateData := bson.M{
		"used_count": voucher.UsedCount,
		"updated_at": voucher.UpdatedAt,
	}
	if err := s.repo.UpdateFields(ctx, filter, bson.M{"$set": updateData}); err != nil {
		return nil, fmt.Errorf("failed to update voucher usage: %w", err)
	}

	return voucher, nil
}

// UpdateVoucher updates one or more fields for a given voucher ID.
func (s *voucherService) UpdateVoucher(ctx context.Context, id primitive.ObjectID, fields map[string]interface{}) error {
	fields["updated_at"] = time.Now()

	filter := bson.M{"_id": id}
	updateData := bson.M{"$set": fields}

	if err := s.repo.UpdateFields(ctx, filter, updateData); err != nil {
		return fmt.Errorf("failed to update voucher: %w", err)
	}
	return nil
}

// GetAllVouchers returns a list of all vouchers.
func (s *voucherService) GetAllVouchers(ctx context.Context) ([]entity.VoucherCode, error) {
	return s.repo.GetAll(ctx)
}
