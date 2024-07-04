package dbrepo

import (
	"context"
	"database/sql"
	"mlvt/internal/models"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) GetVideoByID(id int) (*models.Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := ``

	var video models.Video

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&video.Duration,
		&video.Size,
		//and more
	)

	if err != nil {
		return nil, err
	}

	return &video, nil
}
