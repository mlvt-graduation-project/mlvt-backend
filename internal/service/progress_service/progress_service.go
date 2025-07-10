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
	"mlvt/internal/repo/media_repo"
	"mlvt/internal/repo/progress_repo"
	"mlvt/internal/utility"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgressService interface {
	Create(ctx context.Context, p entity.Progress) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Progress, error)
	GetByFilter(ctx context.Context, qo mongodb.QueryOptions) ([]entity.Progress, error)
	GetProgressByUserID(ctx context.Context, userID uint64, offset int, limit int, searchKey string, progressType []entity.ProgressType, progressStatus []entity.StatusEntity) ([]entity.Progress, int, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, newStatus entity.StatusEntity) error
	UpdateTitle(ctx context.Context, id primitive.ObjectID, newTitle string) error
	DeleteProgress(ctx context.Context, id primitive.ObjectID) error
	UpdateFieldId(ctx context.Context, id primitive.ObjectID, fieldName string, value uint64) error
	GetCountProgressTypeByUserId(ctx context.Context, userID uint64, progressType entity.ProgressType) (int, error)
	GetProgressThumbnails(progresses []entity.Progress) ([]response.ProgressResponse, error)
}

type progressService struct {
	repo      progress_repo.ProgressRepository
	mediaRepo media_repo.MediaRepository
	s3Client  aws.S3ClientInterface
}

func NewProgressService(
	repo progress_repo.ProgressRepository,
	mediaRepo media_repo.MediaRepository,
	s3Client aws.S3ClientInterface,
) ProgressService {
	return &progressService{
		repo:      repo,
		mediaRepo: mediaRepo,
		s3Client:  s3Client,
	}
}

func (s *progressService) Create(ctx context.Context, p entity.Progress) (primitive.ObjectID, error) {
	count, err := s.GetCountProgressTypeByUserId(ctx, p.UserID, p.ProgressType)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to get total count of progress type")
	}
	p.Title = utility.GetProgressTitle(count+1, p.ProgressType)
	return s.repo.Insert(ctx, p)
}

func (s *progressService) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Progress, error) {
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

func (s *progressService) UpdateTitle(ctx context.Context, id primitive.ObjectID, newTitle string) error {
	filter := bson.M{"_id": id}

	updateData := bson.M{
		"title":      newTitle,
		"updated_at": time.Now(),
	}

	return s.repo.UpdateFields(ctx, filter, updateData)
}

func (s *progressService) DeleteProgress(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	updateData := bson.M{
		"is_deleted": true,
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
	offset int,
	limit int,
	searchKey string,
	progressType []entity.ProgressType,
	progressStatus []entity.StatusEntity,
) (
	[]entity.Progress,
	int,
	error,
) {
	filters := []mongodb.FilterCondition{
		{
			Key:       "user_id",
			Operation: mongodb.OpEqual,
			Value:     userID,
		},
	}

	if searchKey != "" {
		filters = append(filters, mongodb.FilterCondition{
			Key:       "title",
			Operation: mongodb.OpRegex,
			Value: primitive.Regex{
				Pattern: ".*" + regexp.QuoteMeta(searchKey) + ".*",
				Options: "i",
			},
		})
	}

	// Nếu có lọc theo progressType thì thêm điều kiện
	if len(progressType) > 0 {
		values := make([]interface{}, len(progressType))
		for i, pt := range progressType {
			values[i] = pt
		}
		filters = append(filters, mongodb.FilterCondition{
			Key:       "progress_type",
			Operation: mongodb.OpIn,
			Value:     values,
		})
	}

	if len(progressStatus) > 0 {
		values := make([]interface{}, len(progressStatus))
		for i, pt := range progressStatus {
			values[i] = pt
		}
		filters = append(filters, mongodb.FilterCondition{
			Key:       "status",
			Operation: mongodb.OpIn,
			Value:     values,
		})
	}

	filters = append(filters, mongodb.FilterCondition{
		Key:       "is_deleted",
		Operation: mongodb.OpNotEqual,
		Value:     true,
	})

	qo := mongodb.QueryOptions{
		Filters: filters,
		Sorts: []mongodb.SortCondition{
			{
				Field:     "created_at",
				Direction: mongodb.SortDesc,
			},
		},
		Limit:  &limit,
		Offset: &offset,
	}

	progressList, err := s.repo.GetByFilter(ctx, qo)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := s.repo.CountByFilter(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count progress: %w", err)
	}

	return progressList, totalCount, nil
}

func (s *progressService) GetProgressThumbnails(
	progresses []entity.Progress,
) (
	[]response.ProgressResponse,
	error,
) {
	var result []response.ProgressResponse
	thumbnailURL := ""

	for _, progress := range progresses {
		if progress.ProgressType == "stt" || progress.ProgressType == "ls" || progress.ProgressType == "fp" {
			video, err := s.mediaRepo.GetVideoByID(progress.OriginalVideoID)
			if err != nil {
				log.Errorf("Failed to get video with ID %d: %v", progress.OriginalVideoID, err)
				return nil, fmt.Errorf("failed to get video with ID %d: %w", progress.OriginalVideoID, err)
			}

			// Prevent nil pointer dereference
			if video == nil {
				log.Errorf("Video with ID %d is nil", progress.OriginalVideoID)
				return nil, fmt.Errorf("video with ID %d not found", progress.OriginalVideoID)
			}

			log.Infof("Video found: Folder=%s, Image=%s", video.Folder, video.Image)

			thumbnailURL, err = s.s3Client.GeneratePresignedDownloadURL("video_frames", video.Image, "image/jpeg")
			if err != nil {
				log.Errorf("Failed to generate presigned download URL for video ID %d: %v", progress.OriginalVideoID, err)
				thumbnailURL = "" // Ensure the response still contains valid data even if thumbnail generation fails
			}
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
			Title:                     progress.Title,
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

func (s *progressService) GetCountProgressTypeByUserId(ctx context.Context, userID uint64, progressType entity.ProgressType) (int, error) {
	filters := []mongodb.FilterCondition{
		{
			Key:       "user_id",
			Operation: mongodb.OpEqual,
			Value:     userID,
		},
		{
			Key:       "progress_type",
			Operation: mongodb.OpEqual,
			Value:     progressType,
		},
	}
	return s.repo.CountByFilter(ctx, filters)
}
