package voucher_repo

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VoucherRepository interface {
	Insert(ctx context.Context, voucher entity.VoucherCode) (primitive.ObjectID, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.VoucherCode, error)
	FindByCode(ctx context.Context, code string) (*entity.VoucherCode, error)
	UpdateVoucher(ctx context.Context, filter interface{}, updatedFields interface{}) error
	GetAll(ctx context.Context) ([]entity.VoucherCode, error)
}

type voucherRepo struct {
	adapter *mongodb.MongoDBAdapter[entity.VoucherCode]
}

func NewVoucherRepo(db *mongodb.MongoDBClient) VoucherRepository {
	return &voucherRepo{
		adapter: mongodb.NewMongoDBAdapter[entity.VoucherCode](
			db.GetClient(),
			"mlvt",
			"voucher_codes",
		),
	}
}

func (r *voucherRepo) Insert(ctx context.Context, voucher entity.VoucherCode) (primitive.ObjectID, error) {
	insertedID, err := r.adapter.InsertOne(voucher)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert voucher: %w", err)
	}
	return insertedID, nil
}

func (r *voucherRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.VoucherCode, error) {
	filter := bson.M{"_id": id}

	result, err := r.adapter.FindOne(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find voucher by ID: %w", err)
	}
	return result, nil
}

func (r *voucherRepo) FindByCode(ctx context.Context, code string) (*entity.VoucherCode, error) {
	filter := bson.M{"code": code}

	result, err := r.adapter.FindOne(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find voucher by code: %w", err)
	}
	return result, nil
}

func (r *voucherRepo) UpdateVoucher(
	ctx context.Context,
	filter interface{},
	updatedFields interface{},
) error {
	if err := r.adapter.UpdateOne(filter, updatedFields); err != nil {
		return fmt.Errorf("failed to update voucher: %w", err)
	}
	return nil
}

func (r *voucherRepo) GetAll(ctx context.Context) ([]entity.VoucherCode, error) {
	// Retrieve all vouchers without a filter
	vouchers, err := r.adapter.Find(bson.M{}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get all vouchers: %w", err)
	}
	return vouchers, nil
}
