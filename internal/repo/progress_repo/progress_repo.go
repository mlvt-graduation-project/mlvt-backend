package progress_repo

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgressRepository interface {
	Insert(ctx context.Context, progress entity.Progress) (primitive.ObjectID, error)
	Get(ctx context.Context, id primitive.ObjectID) (*entity.Progress, error)
	GetByFilter(ctx context.Context, queryOpts mongodb.QueryOptions) ([]entity.Progress, error)
	UpdateFields(ctx context.Context, filter interface{}, updateFields interface{}) error
	CountByFilter(ctx context.Context, filters []mongodb.FilterCondition) (int, error)
}

type progressRepo struct {
	adapter *mongodb.MongoDBAdapter[entity.Progress]
}

func NewProgressRepo(db *mongodb.MongoDBClient) ProgressRepository {
	return &progressRepo{
		adapter: mongodb.NewMongoDBAdapter[entity.Progress](
			db.GetClient(),
			"mlvt",
			"progress",
		),
	}
}

func (r *progressRepo) Insert(ctx context.Context, progress entity.Progress) (primitive.ObjectID, error) {
	insertedID, err := r.adapter.InsertOne(progress)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert progress: %w", err)
	}
	return insertedID, nil
}

func (r *progressRepo) Get(ctx context.Context, id primitive.ObjectID) (*entity.Progress, error) {
	filter := bson.M{"_id": id}

	result, err := r.adapter.FindOne(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find progress: %w", err)
	}
	// result == nil means no document found
	return result, nil
}

func (r *progressRepo) GetByFilter(
	ctx context.Context,
	queryOpts mongodb.QueryOptions,
) ([]entity.Progress, error) {
	filter, findOpts, err := mongodb.BuildQuery(queryOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to build query for progress: %w", err)
	}

	docs, err := r.adapter.Find(filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to query progress: %w", err)
	}
	return docs, nil
}

func (r *progressRepo) UpdateFields(
	ctx context.Context,
	filter interface{},
	updateFields interface{},
) error {
	return r.adapter.UpdateOne(filter, updateFields)
}

func (r *progressRepo) CountByFilter(
	ctx context.Context,
	filters []mongodb.FilterCondition,
) (int, error) {
	// Tạo filter BSON từ FilterCondition
	queryOpts := mongodb.QueryOptions{
		Filters: filters,
	}

	filter, _, err := mongodb.BuildQuery(queryOpts)
	if err != nil {
		return 0, fmt.Errorf("failed to build query for count: %w", err)
	}

	count, err := r.adapter.CountDocuments(filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count progress: %w", err)
	}

	return int(count), nil
}
