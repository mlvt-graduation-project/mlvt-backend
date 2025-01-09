package progress_service

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/repo/progress_repo"
)

type ProgressService interface {
	Create(ctx context.Context, p entity.Progress) (uint64, error)
	GetByID(ctx context.Context, id uint64) (*entity.Progress, error)
	GetByFilter(ctx context.Context, qo mongodb.QueryOptions) ([]entity.Progress, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) error
}

type progressService struct {
	repo progress_repo.ProgressRepository
}

func NewProgressService(repo progress_repo.ProgressRepository) ProgressService {
	return &progressService{
		repo: repo,
	}
}

func (s *progressService) Create(ctx context.Context, p entity.Progress) (uint64, error) {
	return s.repo.Insert(ctx, p)
}

func (s *progressService) GetByID(ctx context.Context, id uint64) (*entity.Progress, error) {
	return s.repo.Get(ctx, id)
}

func (s *progressService) GetByFilter(ctx context.Context, qo mongodb.QueryOptions) ([]entity.Progress, error) {
	return s.repo.GetByFilter(ctx, qo)
}

func (s *progressService) UpdateOne(ctx context.Context, filter interface{}, update interface{}) error {
	return s.repo.UpdateOne(ctx, filter, update)
}
