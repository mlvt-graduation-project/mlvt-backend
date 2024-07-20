package database

import (
	"database/sql"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

//const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}
