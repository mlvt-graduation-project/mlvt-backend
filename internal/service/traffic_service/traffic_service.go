package traffic_service

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/repo/traffic_repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrafficService interface {
	CreateTraffic(ctx context.Context, p entity.Traffic) (primitive.ObjectID, error)
	GetTrafficByID(ctx context.Context, id uint64) (*entity.Traffic, error)
	GetTrafficsByUserID(ctx context.Context, userID uint64) ([]entity.Traffic, error)
}

type trafficService struct {
	repo     traffic_repo.TrafficRepository
	s3Client aws.S3ClientInterface
}

func NewTrafficService(
	repo traffic_repo.TrafficRepository,
	s3Client aws.S3ClientInterface,
) TrafficService {
	return &trafficService{
		repo:     repo,
		s3Client: s3Client,
	}
}

func (s *trafficService) CreateTraffic(ctx context.Context, p entity.Traffic) (primitive.ObjectID, error) {
	return s.repo.InsertTraffic(ctx, p)
}

func (s *trafficService) GetTrafficByID(ctx context.Context, id uint64) (*entity.Traffic, error) {
	return s.repo.GetTraffic(ctx, id)
}

func (s *trafficService) GetTrafficsByUserID(
	ctx context.Context,
	userID uint64,
) (
	[]entity.Traffic,
	error,
) {
	qo := mongodb.QueryOptions{
		Filters: []mongodb.FilterCondition{
			{
				Key:       "user_id",
				Operation: mongodb.OpEqual,
				Value:     userID,
			},
		},
		Sorts: []mongodb.SortCondition{
			{
				Field:     "created_at",
				Direction: mongodb.SortDesc,
			},
		},
		// Not set the Field to return all columns
	}

	return s.repo.GetTrafficsByFilter(ctx, qo)
}
