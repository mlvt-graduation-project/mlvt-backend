package progress_service

import (
	"context"
	"errors"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/pkg/response"
	"mlvt/internal/repo/progress_repo"
	"mlvt/internal/repo/video_repo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgressService interface {
	Create(ctx context.Context, p entity.Progress) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id uint64) (*entity.Progress, error)
	GetByFilter(ctx context.Context, qo mongodb.QueryOptions) ([]entity.Progress, error)
	GetProgressByUserID(ctx context.Context, userID uint64) ([]entity.Progress, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, newStatus entity.StatusEntity) error
	UpdateFieldId(ctx context.Context, id primitive.ObjectID, fieldName string, value uint64) error
	GetProgressThumbnails(progresses []entity.Progress) ([]response.ProgressResponse, error)
}

type progressService struct {
	repo      progress_repo.ProgressRepository
	videoRepo video_repo.VideoRepository
	s3Client  aws.S3ClientInterface
}

func NewProgressService(
	repo progress_repo.ProgressRepository,
	videoRepo video_repo.VideoRepository,
	s3Client aws.S3ClientInterface,
) ProgressService {
	return &progressService{
		repo:      repo,
		videoRepo: videoRepo,
		s3Client:  s3Client,
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

func (s *progressService) GetProgressByUserID(
	ctx context.Context,
	userID uint64,
) (
	[]entity.Progress,
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

	return s.repo.GetByFilter(ctx, qo)
}

func (s *progressService) GetProgressThumbnails(
	progresses []entity.Progress,
) (
	[]response.ProgressResponse,
	error,
) {
	var result []response.ProgressResponse
	for _, progress := range progresses {
		video, err := s.videoRepo.GetVideoByID(progress.OriginalVideoID)
		if err != nil {
			log.Error("Failed to get video id, GetProgressthumbnails funtion")
			return nil, fmt.Errorf("failed to get video with ID %d: %w", progress.OriginalVideoID, err)
		}

		thumbnailURL, err := s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.Image, "image/jpeg")
		if err != nil {
			log.Errorf("Failed to generate presigned download url for video id: %d", progress.OriginalVideoID)
		}

		resp := response.ProgressResponse{
			ID:                        progress.ID,
			UserID:                    progress.UserID,
			ProgressType:              progress.ProgressType,
			OriginalVideoID:           progress.OriginalVideoID,
			OriginalTranscriptionID:   progress.OriginalTranscriptionID,
			TranslatedTranscriptionID: progress.TranslatedTranscriptionID,
			AudioID:                   progress.AudioID,
			ProgressedVideoID:         progress.ProgressedVideoID,
			Status:                    progress.Status,
			CreatedAt:                 progress.CreatedAt,
			UpdatedAt:                 progress.UpdatedAt,
			ThumbnailUrl:              thumbnailURL,
		}

		result = append(result, resp)
	}

	return result, nil
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
