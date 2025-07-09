package ping_repo

import (
	"database/sql"
	"fmt" // Import fmt for Sprintf
	"mlvt/internal/pkg/response"

	"github.com/jmoiron/sqlx"
)

type PingRepository interface {
	PingAudio(id uint64) (*response.PingStatusResponse, error)
	PingVideo(id uint64) (*response.PingStatusResponse, error)
	PingTranscription(id uint64) (*response.PingStatusResponse, error)
}

type pingRepo struct {
	db *sqlx.DB
}

func NewPingRepo(db *sqlx.DB) PingRepository {
	return &pingRepo{db: db}
}

const (
	tableAudios         = "audios"
	tableVideos         = "videos"
	tableTranscriptions = "transcriptions"
)

func (r *pingRepo) pingStatus(table string, id uint64) (*response.PingStatusResponse, error) {
	query := fmt.Sprintf(`SELECT status FROM %s WHERE id = ?`, table)
	query = sqlx.Rebind(sqlx.DOLLAR, query)

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
