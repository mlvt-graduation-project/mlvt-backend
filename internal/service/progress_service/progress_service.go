package progress_service

import (
	"context"
	"errors"
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
	UpdateFieldId(ctx context.Context, id primitive.ObjectID, fieldName string, value uint64) error
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

func (s *progressService) UpdateFieldId(ctx context.Context, id primitive.ObjectID, fieldName string, value uint64) error {
	bsonFieldName, err := isValidProgressIDField(fieldName)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	updateData := bson.M{
		bsonFieldName: value,
		"updated_at":  time.Now(),
	}

	return s.repo.UpdateFields(ctx, filter, updateData)
}

func isValidProgressIDField(fieldName string) (string, error) {
	switch fieldName {
	case "OriginalVideoID":
		return "original_video_id", nil
	case "OriginalTranscriptionID":
		return "original_transcription_id", nil
	case "TranslatedTranscriptionID":
		return "translated_transcription_id", nil
	case "AudioID":
		return "audio_id", nil
	case "ProgressedVideoID":
		return "progressed_video_id", nil
	default:
		return "", errors.New("invalid field name")
	}
}
