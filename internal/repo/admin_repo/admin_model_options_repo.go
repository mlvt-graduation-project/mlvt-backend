package admin_repo

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func (r *adminRepo) GetModelOptionByID(
	ctx context.Context,
	id primitive.ObjectID,
) (
	*entity.ModelOption,
	error,
) {
	filter := bson.M{
		"_id": id,
	}
	result, err := r.modelOptionAdapter.FindOne(filter)
	if err != nil {
		return nil, fmt.Errorf("cannot find document with id %s: %w", id, err)
	}

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return result, nil
}

func (r *adminRepo) UpdateModelOption(
	ctx context.Context,
	filter interface{},
	updatedFields interface{},
) error {
	return r.modelOptionAdapter.UpdateOne(filter, updatedFields)
}
