package service

import (
	"mlvt/internal/entity"

	"github.com/stretchr/testify/mock"
)

// MockVideoService is a mock implementation of the VideoService interface
type MockVideoService struct {
	mock.Mock
}

func (m *MockVideoService) CreateVideo(video *entity.Video) error {
	args := m.Called(video)
	return args.Error(0)
}

func (m *MockVideoService) GetVideoByID(videoID uint64) (*entity.Video, string, string, error) {
	args := m.Called(videoID)
	video, _ := args.Get(0).(*entity.Video)
	return video, args.String(1), args.String(2), args.Error(3)
}

func (m *MockVideoService) ListVideosByUserID(userID uint64) ([]entity.Video, []entity.Frame, error) {
	args := m.Called(userID)
	return args.Get(0).([]entity.Video), args.Get(1).([]entity.Frame), args.Error(2)
}

func (m *MockVideoService) DeleteVideo(videoID uint64) error {
	args := m.Called(videoID)
	return args.Error(0)
}

func (m *MockVideoService) UpdateVideo(video *entity.Video) error {
	args := m.Called(video)
	return args.Error(0)
}

func (m *MockVideoService) UpdateVideoStatus(videoID uint64, status entity.VideoStatus) error {
	args := m.Called(videoID, status)
	return args.Error(0)
}

func (m *MockVideoService) GetVideoStatus(videoID uint64) (entity.VideoStatus, error) {
	args := m.Called(videoID)
	return args.Get(0).(entity.VideoStatus), args.Error(1)
}

func (m *MockVideoService) GeneratePresignedUploadURLForVideo(folder, fileName, fileType string) (string, error) {
	args := m.Called(folder, fileName, fileType)
	return args.String(0), args.Error(1)
}

func (m *MockVideoService) GeneratePresignedUploadURLForImage(folder, fileName, fileType string) (string, error) {
	args := m.Called(folder, fileName, fileType)
	return args.String(0), args.Error(1)
}

func (m *MockVideoService) GeneratePresignedDownloadURLForVideo(videoID uint64) (string, error) {
	args := m.Called(videoID)
	return args.String(0), args.Error(1)
}

func (m *MockVideoService) GeneratePresignedDownloadURLForImage(videoID uint64) (string, error) {
	args := m.Called(videoID)
	return args.String(0), args.Error(1)
}
