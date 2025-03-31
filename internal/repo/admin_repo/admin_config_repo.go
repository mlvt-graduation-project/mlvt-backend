package admin_repo

import (
	"context"
	"errors"
	"mlvt/internal/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

func (r *adminRepo) GetAdminConfig(ctx context.Context) (*entity.AdminConfig, error) {
	result, err := r.adminConfigAdapter.FindOne(nil)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New("cannot find admin config")
	}

	return result, nil
}

func (r *adminRepo) UpdateConfig(
	ctx context.Context,
	filter interface{},
	updateFields interface{},
) error {
	return r.adminConfigAdapter.UpdateOne(filter, updateFields)
}
