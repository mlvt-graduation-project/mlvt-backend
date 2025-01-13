package progress_service

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/repo/progress_repo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgressService interface {
	Create(ctx context.Context, p entity.Progress) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id uint64) (*entity.Progress, error)
	GetByFilter(ctx context.Context, qo mongodb.QueryOptions) ([]entity.Progress, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, newStatus entity.StatusEntity) error
}

type progressService struct {
	repo progress_repo.ProgressRepository
}

func NewProgressService(repo progress_repo.ProgressRepository) ProgressService {
	return &progressService{
		repo: repo,
	}
}

func (s *progressService) Create(ctx context.Context, p entity.Progress) (primitive.ObjectID, error) {
	return s.repo.Insert(ctx, p)
}

func (s *progressService) GetByID(ctx context.Context, id uint64) (*entity.Progress, error) {
	return s.repo.Get(ctx, id)
}

func (s *progressService) GetByFilter(ctx context.Context, qo mongodb.QueryOptions) ([]entity.Progress, error) {
	return s.repo.GetByFilter(ctx, qo)
}

func (s *progressService) UpdateStatus(ctx context.Context, id primitive.ObjectID, newStatus entity.StatusEntity) error {
	filter := bson.M{"_id": id}

	updateData := bson.M{
		"status":     newStatus,
		"updated_at": time.Now(),
	}

	return s.repo.UpdateFields(ctx, filter, updateData)
}
