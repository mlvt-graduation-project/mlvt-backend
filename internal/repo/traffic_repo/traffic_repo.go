package traffic_repo

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrafficRepository interface {
	InsertTraffic(ctx context.Context, traffic entity.Traffic) (primitive.ObjectID, error)
	GetTraffic(ctx context.Context, id uint64) (*entity.Traffic, error)
	GetTrafficsByFilter(ctx context.Context, queryOpts mongodb.QueryOptions) ([]entity.Traffic, error)
	UpdateTrafficFields(ctx context.Context, filter interface{}, updateFields interface{}) error
}

type trafficRepo struct {
	adapter *mongodb.MongoDBAdapter[entity.Traffic]
}

func NewTrafficRepo(db *mongodb.MongoDBClient) TrafficRepository {
	return &trafficRepo{
		adapter: mongodb.NewMongoDBAdapter[entity.Traffic](
			db.GetClient(),
			"mlvt",
			"traffic",
		),
	}
}

func (r *trafficRepo) InsertTraffic(
	ctx context.Context,
	traffic entity.Traffic,
) (
	primitive.ObjectID,
	error,
) {
	insertedID, err := r.adapter.InsertOne(traffic)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert traffic: %w", err)
	}
	return insertedID, nil
}

func (r *trafficRepo) GetTraffic(
	ctx context.Context,
	id uint64,
) (
	*entity.Traffic,
	error,
) {
	filter := bson.M{"_id": id}

	result, err := r.adapter.FindOne(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find traffic: %w", err)
	}
	return result, nil
}

func (r *trafficRepo) GetTrafficsByFilter(
	ctx context.Context,
	queryOpts mongodb.QueryOptions,
) (
	[]entity.Traffic,
	error,
) {
	filter, findOpts, err := mongodb.BuildQuery(queryOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to build query for traffic: %w", err)
	}

	docs, err := r.adapter.Find(filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to query traffic: %w", err)
	}
	return docs, nil
}

func (r *trafficRepo) UpdateTrafficFields(
	ctx context.Context,
	filter interface{},
	updateFields interface{},
) error {
	return r.adapter.UpdateOne(filter, updateFields)
}
