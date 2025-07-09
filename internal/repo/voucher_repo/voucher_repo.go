package voucher_repo

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/utils"
	"regexp"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VoucherRepository interface {
	Insert(ctx context.Context, voucher entity.VoucherCode) (primitive.ObjectID, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.VoucherCode, error)
	FindByCode(ctx context.Context, code string) (*entity.VoucherCode, error)
	UpdateVoucher(ctx context.Context, filter interface{}, updatedFields interface{}) error
	GetAll(ctx context.Context, status string, sortBy string, sortOrder int, searchField string, searchKey string, offset int, limit int) (vouchers []entity.VoucherCode, totalCount int64, err error)
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

func (r *voucherRepo) GetAll(
	ctx context.Context,
	status string,
	sortBy string,
	sortOrder int,
	searchField string,
	searchKey string,
	offset int,
	limit int,
) ([]entity.VoucherCode, int64, error) {
	filter := bson.M{}
	now := time.Now()

	switch status {
		case "ACTIVE":
			filter["expired_time"] = bson.M{"$gt": now}
			filter["$expr"] = bson.M{
				"$lt": []interface{}{"$used_count", "$max_usage"},
			}
		case "EXPIRED":
			filter["expired_time"] = bson.M{"$lte": now}
		case "USED":
			filter["$expr"] = bson.M{
				"$gte": []interface{}{"$used_count", "$max_usage"},
			}
	}
	if searchField != "" && searchKey != "" {
		if utils.IsInListString(searchField, []string{"used_count", "token", "max_usage"}) {
			key, _ := strconv.Atoi(searchKey)
			filter[searchField] = key
		} else {
			filter[searchField] = bson.M{
				"$regex":   regexp.QuoteMeta(searchKey),
				"$options": "i",
			}
		}
	}

	fmt.Println("SortBy: ", sortBy)
	fmt.Println("SortOrder: ", sortOrder)

	// Step 1: Count total
	totalCount, err := r.adapter.CountDocuments(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vouchers: %w", err)
	}

	// Step 2: Find with pagination
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})
	if offset > 0 {
		findOptions.SetSkip(int64(offset))
	}
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}

	results, err := r.adapter.Find(filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get vouchers: %w", err)
	}

	return results, totalCount, nil
}
