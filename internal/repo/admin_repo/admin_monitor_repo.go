package admin_repo

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"time"
)

// #region media

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

// #endregion

// #region pipeline

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
	tttSucceeded, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeTTT, string(entity.StatusSucceeded))
	if err != nil {
		return pipeline, err
	}
	tttFailed, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeTTT, string(entity.StatusFailed))
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
	sttSucceeded, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeSTT, string(entity.StatusSucceeded))
	if err != nil {
		return pipeline, err
	}
	sttFailed, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeSTT, string(entity.StatusFailed))
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
	lsSucceeded, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeLS, string(entity.StatusSucceeded))
	if err != nil {
		return pipeline, err
	}
	lsFailed, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeLS, string(entity.StatusFailed))
	if err != nil {
		return pipeline, err
	}
	pipeline.LS = entity.MonitorMetric{
		Count:     uint64(lsCount),
		Succeeded: uint64(lsSucceeded),
		Failed:    uint64(lsFailed),
	}

	// Count FP
	fpCount, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeFP, "")
	if err != nil {
		return pipeline, err
	}
	fpSucceeded, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeFP, string(entity.StatusSucceeded))
	if err != nil {
		return pipeline, err
	}
	fpFailed, err := r.countByTypeAndStatus(ctx, entity.ProgressTypeFP, string(entity.StatusFailed))
	if err != nil {
		return pipeline, err
	}
	pipeline.FP = entity.MonitorMetric{
		Count:     uint64(fpCount),
		Succeeded: uint64(fpSucceeded),
		Failed:    uint64(fpFailed),
	}

	pipeline.All = entity.MonitorMetric{
		Count:     uint64(sttCount + tttCount + ttsCount + lsSucceeded),
		Succeeded: uint64(sttSucceeded + tttSucceeded + ttsSucceeded + lsSucceeded),
		Failed:    uint64(sttFailed + tttFailed + ttsFailed + lsFailed),
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

// #endregion

// #region traffic

func (r *adminRepo) GetMonitorTraffic(
	ctx context.Context,
	adminID uint64,
	timeType entity.TimePeriodType,
	baseTime time.Time,
) (entity.MonitorTraffics, error) {
	var (
		response     entity.MonitorTraffics
		segmentCount int
		start        time.Time
		end          time.Time
		labels       []string
		err          error
	)

	now := baseTime

	switch timeType {
	case entity.TimePeriodDay:
		segmentCount = 24
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1)   // add 1 day
		labels = buildDayLabels(start) // ["00:00-01:00","01:00-02:00",...,"23:00-00:00"]
	default:
		return response, fmt.Errorf("invalid timeType: %s", timeType)
	}

	filters := []mongodb.FilterCondition{
		{
			Key:       "timestamp",
			Operation: mongodb.OpGTE,
			Value:     start.Unix(),
		},
		{
			Key:       "timestamp",
			Operation: mongodb.OpLTE,
			Value:     end.Unix(),
		},
	}

	trafficList, err := r.trafficAdapter.FindWithQuery(filters)
	if err != nil {
		return response, fmt.Errorf("failed to query traffic: %w", err)
	}

	// prepare usage counters
	usageCounts := make([]uint64, segmentCount)
	totalSeconds := end.Unix() - start.Unix()
	if totalSeconds <= 0 {
		return response, fmt.Errorf("data range is invalid")
	}

	// Tally each document into the appropriate segment
	for _, trafficItem := range trafficList {
		ts := trafficItem.Timestamp
		if ts < start.Unix() || ts >= end.Unix() {
			continue
		}
		offset := ts - start.Unix()
		frac := float64(offset) / float64(totalSeconds)
		index := int(frac * float64(totalSeconds))
		if index >= segmentCount {
			index = segmentCount - 1
		}
		usageCounts[index]++
	}

	// build final segments
	segments := make([]entity.MonitorTraffic, segmentCount)
	var totalUsage uint64
	for i := 0; i < segmentCount; i++ {
		segments[i] = entity.MonitorTraffic{
			Cell:  labels[i],
			Value: usageCounts[i],
		}
		totalUsage += usageCounts[i]
	}

	response = entity.MonitorTraffics{
		Count:   totalUsage,
		Traffic: segments,
	}

	return response, nil
}

func buildDayLabels(day time.Time) []string {
	labels := make([]string, 24)
	for hour := 0; hour < 24; hour++ {
		next := hour + 1
		labels[hour] = fmt.Sprintf("%02d:00-%02d:00", hour, next%24)
	}
	return labels
}

// #endregion
