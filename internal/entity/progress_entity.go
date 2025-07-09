package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Progress struct {
	ID                        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID                    uint64             `json:"user_id" bson:"user_id"`
	ProgressType              ProgressType       `json:"progress_type" bson:"progress_type"`
	OriginalVideoID           uint64             `json:"original_video_id" bson:"original_video_id"`
	OriginalTranscriptionID   uint64             `json:"original_transcription_id" bson:"original_transcription_id"`
	TranslatedTranscriptionID uint64             `json:"translated_transcription_id" bson:"translated_transcription_id"`
	AudioID                   uint64             `json:"audio_id" bson:"audio_id"`
	ProgressedVideoID         uint64             `json:"progressed_video_id" bson:"progressed_video_id"`
	Status                    StatusEntity       `json:"status" bson:"status"`
	CreatedAt                 time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt                 time.Time          `json:"updated_at" bson:"updated_at"`
	Title					  string			 `json:"title" bson:"title"`
}
