package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminConfig struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	ConfigKey string             `json:"config_key" bson:"config_key"`
	STTModel  string             `json:"stt_model" bson:"stt_model"`
	TTTModel  string             `json:"ttt_model" bson:"ttt_model"`
	TTSModel  string             `json:"tts_model" bson:"tts_model"`
	LSModel   string             `json:"ls_model" bson:"ls_model"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
