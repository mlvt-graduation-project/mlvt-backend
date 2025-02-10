package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelOptionStatus string

const (
	Available   ModelOptionStatus = "available"
	Unavailable ModelOptionStatus = "unavailable"
)

type ModelOption struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	ModelType   ProgressType       `json:"model_type" bson:"model_type"`
	ModelName   string             `json:"model_name" bson:"model_name"`
	Description string             `json:"description" bson:"description"`
	Status      ModelOptionStatus  `json:"status" bson:"status"`
	Token       int64              `json:"token" bson:"token"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
