package repo

import (
	"mlvt/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockVideoRepository struct {
	mock.Mock
}

func (m *MockVideoRepository) CreateVideo(video *entity.Video) error {
	args := m.Called(video)
	return args.Error(0)
}

func (m *MockVideoRepository) GetVideoByID(videoID uint64) (*entity.Video, error) {
	args := m.Called(videoID)
	return args.Get(0).(*entity.Video), args.Error(1)
}

func (m *MockVideoRepository) ListVideosByUserID(userID uint64) ([]entity.Video, error) {
	args := m.Called(userID)
	return args.Get(0).([]entity.Video), args.Error(1)
}

func (m *MockVideoRepository) DeleteVideo(videoID uint64) error {
	args := m.Called(videoID)
	return args.Error(0)
}

func (m *MockVideoRepository) UpdateVideo(video *entity.Video) error {
	args := m.Called(video)
	return args.Error(0)
}

func (m *MockVideoRepository) UpdateVideoStatus(videoID uint64, status entity.VideoStatus) error {
	args := m.Called(videoID, status)
	return args.Error(0)
}

func (m *MockVideoRepository) GetVideoStatus(videoID uint64) (entity.VideoStatus, error) {
	args := m.Called(videoID)
	return args.Get(0).(entity.VideoStatus), args.Error(1)
}
