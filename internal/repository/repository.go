package repository

import (
	"database/sql"
	"mlvt/internal/models"
)

type DatabaseRepository interface {
	Connection() *sql.DB
	GetVideoByID(id int) (*models.Video, error)
}
