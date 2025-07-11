package media_repo

import (
	"context"
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/zap-logging/log"
	"slices"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func (r *mediaRepo) GetAllMedia(
	userID uint64,
	searchKey string,
	limit int,
	offset int,
	status []entity.StatusEntity,
	mediaType []entity.MediaType,
) ([]entity.Video, []entity.Audio, []entity.Transcription, int, error) {
	type media struct {
		ID         int       `db:"id"`
		Type       string    `db:"type"`
		Title      string    `db:"title"`
		CreatedAt  time.Time `db:"created_at"`
		TotalCount int       `db:"total"`
	}
	type queryBatch struct {
		mediaType string
		query     string
		IDs       []int
		Result    interface{}
	}

	var (
		generalMedia     []media
		resVideo         []entity.Video
		resAudio         []entity.Audio
		resTranscription []entity.Transcription
		videoIDs         []int
		audioIDs         []int
		transcriptionIDs []int
	)

	// Search keyword
	likePattern := "%" + searchKey + "%"

	// --- Build subqueries & args ---
	var queries []string
	var args []interface{}

	if slices.Contains(mediaType, entity.MediaTypeVideo) {
		queries = append(queries, `
		SELECT id, 'video' AS type, title, created_at
		FROM videos
		WHERE is_deleted = false AND user_id = ? AND status IN (?) AND title ILIKE ?
	`)
		args = append(args, userID, status, likePattern)
	}
	if slices.Contains(mediaType, entity.MediaTypeAudio) {
		queries = append(queries, `
		SELECT id, 'audio' AS type, title, created_at
		FROM audios
		WHERE is_deleted = false AND user_id = ? AND status IN (?) AND title ILIKE ?
	`)
		args = append(args, userID, status, likePattern)
	}
	if slices.Contains(mediaType, entity.MediaTypeText) {
		queries = append(queries, `
		SELECT id, 'text' AS type, title, created_at
		FROM transcriptions
		WHERE is_deleted = false AND user_id = ? AND status IN (?) AND title ILIKE ?
	`)
		args = append(args, userID, status, likePattern)
	}

	// No allowed media
	if len(queries) == 0 {
		return nil, nil, nil, 0, nil
	}

	// Ghép query tổng
	innerQuery := strings.Join(queries, "\nUNION ALL\n")
	finalQuery := fmt.Sprintf(`
		SELECT id, type, title, created_at, COUNT(*) OVER () AS total
		FROM (
			%s
		) AS all_media
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, innerQuery)

	// Thêm limit, offset vào args
	args = append(args, limit, offset)

	// Sử dụng sqlx.In để mở rộng cho slice trong status (IN (?))
	rawQuery, flatArgs, err := sqlx.In(finalQuery, args...)
	if err != nil {
		return nil, nil, nil, 0, fmt.Errorf("failed to build query: %w", err)
	}
	rawQuery = r.db.Rebind(rawQuery)

	err = r.db.SelectContext(context.Background(), &generalMedia, rawQuery, flatArgs...)
	if err != nil {
		if err == sql.ErrNoRows {
			return resVideo, resAudio, resTranscription, 0, nil
		}
		log.Errorf("error querying general media list: %v", err)
		return nil, nil, nil, 0, err
	}

	// Phân loại kết quả
	for _, item := range generalMedia {
		switch item.Type {
		case string(entity.MediaTypeVideo):
			videoIDs = append(videoIDs, item.ID)
		case string(entity.MediaTypeAudio):
			audioIDs = append(audioIDs, item.ID)
		case string(entity.MediaTypeText):
			transcriptionIDs = append(transcriptionIDs, item.ID)
		}
	}

	total := 0
	if len(generalMedia) > 0 {
		total = generalMedia[0].TotalCount
	}

	// Dùng con trỏ để select vào đúng slice
	temp := []queryBatch{
		{
			mediaType: string(entity.MediaTypeVideo),
			query: `SELECT 
						id,
						COALESCE(original_video_id, 0) AS original_video_id,
						COALESCE(audio_id, 0) AS audio_id,
						COALESCE(duration, 0) AS duration,
						COALESCE(title, '') AS title,
						COALESCE(description, '') AS description,
						COALESCE(file_name, '') AS file_name,
						COALESCE(folder, '') AS folder,
						COALESCE(image, '') AS image,
						status,
						user_id,
						created_at,
						updated_at,
						is_deleted
					FROM videos
					WHERE id IN (?)
			`,
			IDs:    videoIDs,
			Result: &resVideo,
		},
		{
			mediaType: string(entity.MediaTypeAudio),
			query: `SELECT 
						id,
						COALESCE(title, '') AS title,
						COALESCE(file_name, '') AS file_name,
						COALESCE(folder, '') AS folder,
						COALESCE(transcription_id, 0) AS transcription_id,
						status,
						user_id,
						created_at,
						updated_at,
						is_deleted
					FROM audios
					WHERE id IN (?)
			`,
			IDs:    audioIDs,
			Result: &resAudio,
		},
		{
			mediaType: string(entity.MediaTypeText),
			query: `SELECT 
						id,
						COALESCE(original_transcription_id, 0) AS original_transcription_id,
						COALESCE(video_id, 0) AS video_id,
						COALESCE(user_id, 0) AS user_id,
						COALESCE(title, '') AS title,
						COALESCE(lang, '') AS lang,
						COALESCE(folder, '') AS folder,
						COALESCE(file_name, '') AS file_name,
						status,
						created_at,
						updated_at
					FROM transcriptions
					WHERE id IN (?)
			`,
			IDs:    transcriptionIDs,
			Result: &resTranscription,
		},
	}

	for _, batch := range temp {
		if len(batch.IDs) == 0 {
			continue
		}
		subQuery, subArgs, err := sqlx.In(batch.query, batch.IDs)
		if err != nil {
			log.Errorf("error building IN query for %s: %v", batch.mediaType, err)
			continue
		}
		subQuery = r.db.Rebind(subQuery)
		if err := r.db.SelectContext(context.Background(), batch.Result, subQuery, subArgs...); err != nil {
			log.Errorf("error selecting %s items: %v", batch.mediaType, err)
		}
	}

	return resVideo, resAudio, resTranscription, total, nil
}
