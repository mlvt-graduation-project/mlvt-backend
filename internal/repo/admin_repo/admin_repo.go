package admin_repo

import (
	"context"
	"errors"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository interface {
	GetAdminConfig(ctx context.Context) (*entity.AdminConfig, error)
	UpdateConfig(
		ctx context.Context,
		filter interface{},
		updateFields interface{},
	) error
	AddModelOptions(
		ctx context.Context,
		modelOption entity.ModelOption,
	) (
		primitive.ObjectID,
		error,
	)
	LoadModelOptions(
		ctx context.Context,
		queryOpts mongodb.QueryOptions,
	) (
		[]entity.ModelOption,
		error,
	)
	UpdateModelOption(
		ctx context.Context,
		filter interface{},
		updatedFields interface{},
	) error
}

type adminRepo struct {
	adminConfigAdapter *mongodb.MongoDBAdapter[entity.AdminConfig]
	modelOptionAdapter *mongodb.MongoDBAdapter[entity.ModelOption]
}

func NewAminRepo(db *mongodb.MongoDBClient) AdminRepository {
	return &adminRepo{
		adminConfigAdapter: mongodb.NewMongoDBAdapter[entity.AdminConfig](
			db.GetClient(),
			"mlvt",
			"admin_config",
		),
		modelOptionAdapter: mongodb.NewMongoDBAdapter[entity.ModelOption](
			db.GetClient(),
			"mlvt",
			"model_option",
		),
	}
}

func (r *adminRepo) GetAdminConfig(ctx context.Context) (*entity.AdminConfig, error) {
	result, err := r.adminConfigAdapter.FindOne(nil)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("cannot find admin config")
	}

	return result, nil
}

func (r *adminRepo) UpdateConfig(
	ctx context.Context,
	filter interface{},
	updateFields interface{},
) error {
	return r.adminConfigAdapter.UpdateOne(filter, updateFields)
}

func (r *adminRepo) AddModelOptions(
	ctx context.Context,
	modelOption entity.ModelOption,
) (
	primitive.ObjectID,
	error,
) {
	insertedID, err := r.modelOptionAdapter.InsertOne(modelOption)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert model option: %w", err)
	}
	return insertedID, nil
}

func (r *adminRepo) LoadModelOptions(
	ctx context.Context,
	queryOpts mongodb.QueryOptions,
) (
	[]entity.ModelOption,
	error,
) {
	filter, findOpts, err := mongodb.BuildQuery(queryOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to build query for progress: %w", err)
	}

	docs, err := r.modelOptionAdapter.Find(filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to query progress: %w", err)
	}

	return docs, nil
}

func (r *adminRepo) UpdateModelOption(
	ctx context.Context,
	filter interface{},
	updatedFields interface{},
) error {
	return r.modelOptionAdapter.UpdateOne(filter, updatedFields)
}
