package service

import (
	"context"
	"mlvt/internal/pkg/response"
	"mlvt/internal/repo"
)

type PingService interface {
	PingSpeechToText(ctx context.Context, id uint64) (*response.PingStatusResponse, error)
	PingTextToText(ctx context.Context, id uint64) (*response.PingStatusResponse, error)
	PingTextToSpeech(ctx context.Context, id uint64) (*response.PingStatusResponse, error)
	PingVoiceCloning(ctx context.Context, id uint64) (*response.PingStatusResponse, error)
	PingLipSync(ctx context.Context, id uint64) (*response.PingStatusResponse, error)
	PingFullPipeline(ctx context.Context, id uint64) (*response.PingStatusResponse, error)
}

type pingService struct {
	repository repo.PingRepository
}

func NewPingService(repository repo.PingRepository) PingService {
	return &pingService{
		repository: repository,
	}
}

func (s *pingService) PingSpeechToText(ctx context.Context, id uint64) (*response.PingStatusResponse, error) {
	return s.repository.PingTranscription(id)
}

func (s *pingService) PingTextToText(ctx context.Context, id uint64) (*response.PingStatusResponse, error) {
	return s.repository.PingTranscription(id)
}

func (s *pingService) PingTextToSpeech(ctx context.Context, id uint64) (*response.PingStatusResponse, error) {
	return s.repository.PingAudio(id)
}

func (s *pingService) PingVoiceCloning(ctx context.Context, id uint64) (*response.PingStatusResponse, error) {
	return s.repository.PingAudio(id)
}

func (s *pingService) PingLipSync(ctx context.Context, id uint64) (*response.PingStatusResponse, error) {
	return s.repository.PingVideo(id)
}

func (s *pingService) PingFullPipeline(ctx context.Context, id uint64) (*response.PingStatusResponse, error) {
	return s.repository.PingVideo(id)
}
