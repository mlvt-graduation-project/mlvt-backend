package service

import (
	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/repo"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestRepoAndS3Client() (*repo.MockVideoRepository, *aws.MockS3Client) {
	repo := new(repo.MockVideoRepository)
	s3Client := new(aws.MockS3Client)
	return repo, s3Client
}

func TestCreateVideoService(t *testing.T) {
	videoRepo, s3Client := setupTestRepoAndS3Client()
	videoService := NewVideoService(videoRepo, s3Client)

	video := &entity.Video{
		Title:       "Test Video",
		Duration:    120,
		Description: "Test Description",
		FileName:    "test.mp4",
		Folder:      "test_folder",
		Image:       "test_image.jpg",
		Status:      entity.StatusRaw,
		UserID:      1,
	}

	videoRepo.On("CreateVideo", video).Return(nil)
	err := videoService.CreateVideo(video)
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
}

func TestGetVideoByIDService(t *testing.T) {
	videoRepo, s3Client := setupTestRepoAndS3Client()
	videoService := NewVideoService(videoRepo, s3Client)

	video := &entity.Video{
		ID:          1,
		Title:       "Test Video",
		Duration:    120,
		Description: "Test Description",
		FileName:    "test.mp4",
		Folder:      "test_folder",
		Image:       "test_image.jpg",
		Status:      entity.StatusRaw,
		UserID:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	videoRepo.On("GetVideoByID", uint64(1)).Return(video, nil)
	s3Client.On("GeneratePresignedURL", video.Folder, video.FileName, "video/mp4").Return("https://s3.amazonaws.com/test_video.mp4", nil)
	s3Client.On("GeneratePresignedURL", video.Folder, video.Image, "image/jpeg").Return("https://s3.amazonaws.com/test_image.jpg", nil)

	result, videoURL, imageURL, err := videoService.GetVideoByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "https://s3.amazonaws.com/test_video.mp4", videoURL)
	assert.Equal(t, "https://s3.amazonaws.com/test_image.jpg", imageURL)
	videoRepo.AssertExpectations(t)
	s3Client.AssertExpectations(t)
}

func TestListVideosByUserIDService(t *testing.T) {
	videoRepo, s3Client := setupTestRepoAndS3Client()
	videoService := NewVideoService(videoRepo, s3Client)

	video1 := entity.Video{
		ID:          1,
		Title:       "Test Video 1",
		Duration:    120,
		Description: "Test Description 1",
		FileName:    "test1.mp4",
		Folder:      "test_folder_1",
		Image:       "test_image_1.jpg",
		Status:      entity.StatusRaw,
		UserID:      1,
	}
	video2 := entity.Video{
		ID:          2,
		Title:       "Test Video 2",
		Duration:    150,
		Description: "Test Description 2",
		FileName:    "test2.mp4",
		Folder:      "test_folder_2",
		Image:       "test_image_2.jpg",
		Status:      entity.StatusProcessing,
		UserID:      1,
	}

	videos := []entity.Video{video1, video2}
	videoRepo.On("ListVideosByUserID", uint64(1)).Return(videos, nil)
	s3Client.On("GeneratePresignedURL", video1.Folder, video1.Image, "image/jpeg").Return("https://s3.amazonaws.com/test_image_1.jpg", nil)
	s3Client.On("GeneratePresignedURL", video2.Folder, video2.Image, "image/jpeg").Return("https://s3.amazonaws.com/test_image_2.jpg", nil)

	resultVideos, frames, err := videoService.ListVideosByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, resultVideos, 2)
	assert.Len(t, frames, 2)
	assert.Equal(t, "https://s3.amazonaws.com/test_image_1.jpg", frames[0].Link)
	assert.Equal(t, "https://s3.amazonaws.com/test_image_2.jpg", frames[1].Link)
	videoRepo.AssertExpectations(t)
	s3Client.AssertExpectations(t)
}

func TestDeleteVideoService(t *testing.T) {
	videoRepo, s3Client := setupTestRepoAndS3Client()
	videoService := NewVideoService(videoRepo, s3Client)

	videoRepo.On("DeleteVideo", uint64(1)).Return(nil)
	err := videoService.DeleteVideo(1)
	assert.NoError(t, err)
	videoRepo.AssertExpectations(t)
}
