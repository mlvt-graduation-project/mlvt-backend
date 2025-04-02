package admin_repo

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
)

func (r *adminRepo) GetMonitorDataType(
	ctx context.Context,
) (entity.MonitorDataType, error) {
	var result entity.MonitorDataType

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM videos",
	).Scan(&result.Videos.Count); err != nil {
		return result, fmt.Errorf("count videos: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM videos WHERE status = 'succeeded'",
	).Scan(&result.Videos.Succeeded); err != nil {
		return result, fmt.Errorf("count videos succeeded: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM videos WHERE status = 'failed'",
	).Scan(&result.Videos.Failed); err != nil {
		return result, fmt.Errorf("count videos failed: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM audios",
	).Scan(&result.Audios.Count); err != nil {
		return result, fmt.Errorf("count audios: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM audios WHERE status = 'succeeded'",
	).Scan(&result.Audios.Succeeded); err != nil {
		return result, fmt.Errorf("count audios succeeded: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM audios WHERE status = 'failed'",
	).Scan(&result.Audios.Failed); err != nil {
		return result, fmt.Errorf("count audios failed: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM transcriptions",
	).Scan(&result.Texts.Count); err != nil {
		return result, fmt.Errorf("count transcriptions: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM transcriptions WHERE status = 'succeeded'",
	).Scan(&result.Texts.Succeeded); err != nil {
		return result, fmt.Errorf("count transcriptions succeeded: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM transcriptions WHERE status = 'failed'",
	).Scan(&result.Texts.Failed); err != nil {
		return result, fmt.Errorf("count transcriptions failed: %w", err)
	}

	return result, nil
}

// for user usage, will move to public_monitor in future
func (r *adminRepo) GetMonitorDataTypeByUserID(
	ctx context.Context,
	userID uint64,
) (entity.MonitorDataType, error) {
	var result entity.MonitorDataType

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM videos WHERE user_id = ?",
		userID,
	).Scan(&result.Videos.Count); err != nil {
		return result, fmt.Errorf("count videos: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM videos WHERE user_id = ? AND status = 'succeeded'",
		userID,
	).Scan(&result.Videos.Succeeded); err != nil {
		return result, fmt.Errorf("count videos succeeded: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM videos WHERE user_id = ? AND status = 'failed'",
		userID,
	).Scan(&result.Videos.Failed); err != nil {
		return result, fmt.Errorf("count videos failed: %w", err)
	}

	// Audios
	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM audios WHERE user_id = ?",
		userID,
	).Scan(&result.Audios.Count); err != nil {
		return result, fmt.Errorf("count audios: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM audios WHERE user_id = ? AND status = 'succeeded'",
		userID,
	).Scan(&result.Audios.Succeeded); err != nil {
		return result, fmt.Errorf("count audios succeeded: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM audios WHERE user_id = ? AND status = 'failed'",
		userID,
	).Scan(&result.Audios.Failed); err != nil {
		return result, fmt.Errorf("count audios failed: %w", err)
	}

	// Transcriptions (text)
	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM transcriptions WHERE user_id = ?",
		userID,
	).Scan(&result.Texts.Count); err != nil {
		return result, fmt.Errorf("count transcriptions: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM transcriptions WHERE user_id = ? AND status = 'succeeded'",
		userID,
	).Scan(&result.Texts.Succeeded); err != nil {
		return result, fmt.Errorf("count transcriptions succeeded: %w", err)
	}

	if err := r.dbSqlite.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM transcriptions WHERE user_id = ? AND status = 'failed'",
		userID,
	).Scan(&result.Texts.Failed); err != nil {
		return result, fmt.Errorf("count transcriptions failed: %w", err)
	}

	return result, nil
}

func (r *adminRepo) GetMonitorPipeline(ctx context.Context) (entity.MonitorPipeline, error) {
	var pipeline entity.MonitorPipeline

	// count TTS
	ttsCount, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeTTS, "")
	if err != nil {
		return pipeline, err
	}
	ttsSucceeded, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeTTS, string(entity.StatusSucceeded))
	if err != nil {
		return pipeline, err
	}
	ttsFailed, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeTTS, string(entity.StatusFailed))
	if err != nil {
		return pipeline, err
	}
	pipeline.TTS = entity.MonitorMetric{
		Count:     uint64(ttsCount),
		Succeeded: uint64(ttsSucceeded),
		Failed:    uint64(ttsFailed),
	}
	return pipeline, nil
}

func (r *adminRepo) countByTypeAndStatus(ctx context.Context, progressType entity.ProgressType, status string) (int, error) {
	filters := []mongodb.FilterCondition{
		{
			Key:       "progress_type",
			Operation: mongodb.OpEqual,
			Value:     progressType,
		},
	}
	if status != "" {
		filters = append(filters, mongodb.FilterCondition{
			Key:       "status",
			Operation: mongodb.OpEqual,
			Value:     status,
		})
	}

	results, err := r.progressAdapter.FindWithQuery(filters)
	if err != nil {
		return 0, fmt.Errorf("FindWithQuery error: %w", err)
	}

	return len(results), nil
}
