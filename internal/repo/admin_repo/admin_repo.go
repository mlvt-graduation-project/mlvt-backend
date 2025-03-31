package admin_repo

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminRepository interface {
	GetAdminConfig(ctx context.Context) (*entity.AdminConfig, error)
	UpdateConfig(
		ctx context.Context,
		filter interface{},
		updateFields interface{},
	) error
	AddModelOptions(
		ctx context.Context,
		modelOption entity.ModelOption,
	) (
		primitive.ObjectID,
		error,
	)
	LoadModelOptions(
		ctx context.Context,
		queryOpts mongodb.QueryOptions,
	) (
		[]entity.ModelOption,
		error,
	)
	UpdateModelOption(
		ctx context.Context,
		filter interface{},
		updatedFields interface{},
	) error
	GetModelOptionByID(
		ctx context.Context,
		id primitive.ObjectID,
	) (
		*entity.ModelOption,
		error,
	)
}

type adminRepo struct {
	adminConfigAdapter *mongodb.MongoDBAdapter[entity.AdminConfig]
	modelOptionAdapter *mongodb.MongoDBAdapter[entity.ModelOption]
}

func NewAminRepo(db *mongodb.MongoDBClient) AdminRepository {
	return &adminRepo{
		adminConfigAdapter: mongodb.NewMongoDBAdapter[entity.AdminConfig](
			db.GetClient(),
			"mlvt",
			"admin_config",
		),
		modelOptionAdapter: mongodb.NewMongoDBAdapter[entity.ModelOption](
			db.GetClient(),
			"mlvt",
			"model_option",
		),
	}
}
