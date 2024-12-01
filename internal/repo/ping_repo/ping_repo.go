package ping_repo

import (
	"database/sql"
	"fmt" // Import fmt for Sprintf
	"mlvt/internal/pkg/response"
)

type PingRepository interface {
	PingAudio(id uint64) (*response.PingStatusResponse, error)
	PingVideo(id uint64) (*response.PingStatusResponse, error)
	PingTranscription(id uint64) (*response.PingStatusResponse, error)
}

type pingRepo struct {
	db *sql.DB
}

func NewPingRepo(db *sql.DB) PingRepository {
	return &pingRepo{db: db}
}

const (
	tableAudios         = "audios"
	tableVideos         = "videos"
	tableTranscriptions = "transcriptions"
)

func (r *pingRepo) pingStatus(table string, id uint64) (*response.PingStatusResponse, error) {
	query := fmt.Sprintf(`SELECT status FROM %s WHERE id = ?`, table)
	row := r.db.QueryRow(query, id)
	statusRes := &response.PingStatusResponse{}
	err := row.Scan(&statusRes.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return &response.PingStatusResponse{}, nil
		}
		return nil, err
	}
	return statusRes, nil
}

func (r *pingRepo) PingAudio(id uint64) (*response.PingStatusResponse, error) {
	return r.pingStatus(tableAudios, id)
}

func (r *pingRepo) PingVideo(id uint64) (*response.PingStatusResponse, error) {
	return r.pingStatus(tableVideos, id)
}

func (r *pingRepo) PingTranscription(id uint64) (*response.PingStatusResponse, error) {
	return r.pingStatus(tableTranscriptions, id)
}
