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

	// Count TTT
	tttCount, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeTTT, "")
	if err != nil {
		return pipeline, err
	}
	tttSucceeded, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeTTT, "succeeded")
	if err != nil {
		return pipeline, err
	}
	tttFailed, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeTTT, "failed")
	if err != nil {
		return pipeline, err
	}
	pipeline.TTT = entity.MonitorMetric{
		Count:     uint64(tttCount),
		Succeeded: uint64(tttSucceeded),
		Failed:    uint64(tttFailed),
	}

	// Count STT
	sttCount, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeSTT, "")
	if err != nil {
		return pipeline, err
	}
	sttSucceeded, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeSTT, "succeeded")
	if err != nil {
		return pipeline, err
	}
	sttFailed, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeSTT, "failed")
	if err != nil {
		return pipeline, err
	}
	pipeline.STT = entity.MonitorMetric{
		Count:     uint64(sttCount),
		Succeeded: uint64(sttSucceeded),
		Failed:    uint64(sttFailed),
	}

	// Count LS
	lsCount, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeLS, "")
	if err != nil {
		return pipeline, err
	}
	lsSucceeded, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeLS, "succeeded")
	if err != nil {
		return pipeline, err
	}
	lsFailed, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeLS, "failed")
	if err != nil {
		return pipeline, err
	}
	pipeline.LS = entity.MonitorMetric{
		Count:     uint64(lsCount),
		Succeeded: uint64(lsSucceeded),
		Failed:    uint64(lsFailed),
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
