package admin_repo

import (
	"context"
	"database/sql"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminRepository interface {
	// Config
	GetAdminConfig(ctx context.Context) (*entity.AdminConfig, error)
	UpdateConfig(ctx context.Context, filter interface{}, updateFields interface{}) error

	// Model options
	AddModelOptions(ctx context.Context, modelOption entity.ModelOption) (primitive.ObjectID, error)
	LoadModelOptions(ctx context.Context, queryOpts mongodb.QueryOptions) ([]entity.ModelOption, error)
	UpdateModelOption(ctx context.Context, filter interface{}, updatedFields interface{}) error
	GetModelOptionByID(ctx context.Context, id primitive.ObjectID) (*entity.ModelOption, error)

	// Monitor
	GetMonitorDataType(ctx context.Context) (entity.MonitorDataType, error)
	GetMonitorPipeline(ctx context.Context) (entity.MonitorPipeline, error)
}

type adminRepo struct {
	// admin config
	adminConfigAdapter *mongodb.MongoDBAdapter[entity.AdminConfig]
	modelOptionAdapter *mongodb.MongoDBAdapter[entity.ModelOption]

	// admin monitor
	progressAdapter *mongodb.MongoDBAdapter[entity.Progress]
	trafficAdapter  *mongodb.MongoDBAdapter[entity.Traffic]
	dbSqlite        *sql.DB
}

func NewAminRepo(dbMongo *mongodb.MongoDBClient, dbSqlite *sql.DB) AdminRepository {
	return &adminRepo{
		adminConfigAdapter: mongodb.NewMongoDBAdapter[entity.AdminConfig](
			dbMongo.GetClient(),
			"mlvt",
			"admin_config",
		),
		modelOptionAdapter: mongodb.NewMongoDBAdapter[entity.ModelOption](
			dbMongo.GetClient(),
			"mlvt",
			"model_option",
		),
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
