package admin_repo

import (
	"database/sql"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
)

type AdminMonitorRepository interface {
}

type adminMonitorRepo struct {
	progressAdapter *mongodb.MongoDBAdapter[entity.Progress]
	trafficAdapter  *mongodb.MongoDBAdapter[entity.Traffic]
	dbSqlite        *sql.DB
}

func NewAdminMonitorRepo(dbMongo *mongodb.MongoDBClient, dbSqlite *sql.DB) AdminMonitorRepository {
	return &adminMonitorRepo{
		dbSqlite: dbSqlite,
		progressAdapter: mongodb.NewMongoDBAdapter[entity.Progress](
			dbMongo.GetClient(),
			"mlvt",
			"progress",
		),
		trafficAdapter: mongodb.NewMongoDBAdapter[entity.Traffic](
			dbMongo.GetClient(),
			"mlvt",
			"traffic",
		),
	}
}
