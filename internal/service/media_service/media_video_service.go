package media_service

import (
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/pkg/response"
	"mlvt/internal/utility"
)

func (s *mediaService) CreateVideo(video *entity.Video, isFullPipeline bool) (uint64, error) {
	count, err := s.mediaRepo.GetCountVideosByUserId(video.UserID)
	if err != nil {
		return 0, fmt.Errorf("error querying total videos of user id")
	}
	video.Title = utility.GetMediaTitle(entity.MediaTypeVideo, isFullPipeline, false, count+1)
	id, err := s.mediaRepo.CreateVideo(video)
	return id, err
}

func (s *mediaService) GetVideoByID(videoID uint64) (*entity.Video, string, string, error) {
	video, err := s.mediaRepo.GetVideoByID(videoID)
	if err != nil {
		return nil, "", "", err
	}
	if video == nil {
		return nil, "", "", fmt.Errorf("video not found")
	}

	// Generate presigned URLs for video and image
	videoURL, err := s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.FileName, "video/mp4")
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate presigned video URL: %v", err)
	}
	imageURL, err := s.s3Client.GeneratePresignedDownloadURL(env.EnvConfig.VideoFramesFolder, video.Image, "image/jpeg")
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate presigned image URL: %v", err)
	}

	return video, videoURL, imageURL, nil
}

func (s *mediaService) ListVideosByUserID(userID uint64) ([]response.ListVideosByUserIDResponse, error) {
	// Fetch the videos for the user
	videos, err := s.mediaRepo.ListVideosByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list videos for user %d: %v", userID, err)
	}

	var videoWithURLsList []response.ListVideosByUserIDResponse
	for _, video := range videos {
		if video.Status != "raw" && video.Status != "succeeded" {
			continue
		}
		// Generate the presigned URL for the video
		videoURL, err := s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.FileName, "video/mp4")
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned video URL for video ID %d: %v", video.ID, err)
		}

		// Generate the presigned URL for the image/frame
		imageURL, err := s.s3Client.GeneratePresignedDownloadURL(env.EnvConfig.VideoFramesFolder, video.Image, "image/jpeg")
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned image URL for video ID %d: %v", video.ID, err)
		}

		// Append to the list
		videoWithURLsList = append(videoWithURLsList, response.ListVideosByUserIDResponse{
			Video:    video,
			VideoURL: videoURL,
			ImageURL: imageURL,
		})
	}

	return videoWithURLsList, nil
}

func (s *mediaService) ListVideosByUserIDAdvance(userID uint64, searchKey string, limit int, offset int, status []entity.StatusEntity) ([]response.ListVideosByUserIDResponse, error) {
	// Fetch the videos for the user
	videos, err := s.mediaRepo.ListVideosByUserIDAdvance(userID, searchKey, limit, offset, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list videos for user %d: %v", userID, err)
	}

	var videoWithURLsList []response.ListVideosByUserIDResponse
	for _, video := range videos {
		if video.Status != "raw" && video.Status != "succeeded" {
			continue
		}
		// Generate the presigned URL for the video
		videoURL, err := s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.FileName, "video/mp4")
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned video URL for video ID %d: %v", video.ID, err)
		}

		// Generate the presigned URL for the image/frame
		imageURL, err := s.s3Client.GeneratePresignedDownloadURL(env.EnvConfig.VideoFramesFolder, video.Image, "image/jpeg")
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned image URL for video ID %d: %v", video.ID, err)
		}

		// Append to the list
		videoWithURLsList = append(videoWithURLsList, response.ListVideosByUserIDResponse{
			Video:    video,
			VideoURL: videoURL,
			ImageURL: imageURL,
		})
	}

	return videoWithURLsList, nil
}

func (s *mediaService) DeleteVideo(videoID uint64) error {
	// Fetch the video record to get the file names
	video, err := s.mediaRepo.GetVideoByID(videoID)
	if err != nil {
		return fmt.Errorf("failed to fetch video: %v", err)
	}
	if video == nil {
		return fmt.Errorf("video not found")
	}

	// Begin deletion process
	// 1. Delete the video and frame files from S3
	err = s.s3Client.DeleteFile(video.Folder, video.FileName)
	if err != nil {
		return fmt.Errorf("failed to delete video file from S3: %v", err)
	}

	err = s.s3Client.DeleteFile(env.EnvConfig.VideoFramesFolder, video.Image)
	if err != nil {
		return fmt.Errorf("failed to delete frame image from S3: %v", err)
	}

	// 2. Delete the video record from the database
	err = s.mediaRepo.DeleteVideo(videoID)
	if err != nil {
		return fmt.Errorf("failed to delete video from database: %v", err)
	}

	return nil
}

func (s *mediaService) UpdateVideo(video *entity.Video) error {
	return s.mediaRepo.UpdateVideo(video)
}

func (s *mediaService) UpdateVideoStatus(videoID uint64, status entity.StatusEntity) error {
	return s.mediaRepo.UpdateVideoStatus(videoID, status)
}
func (s *mediaService) GetVideoStatus(videoID uint64) (entity.StatusEntity, error) {
	return s.mediaRepo.GetVideoStatus(videoID)
}

// GeneratePresignedUploadURLForVideo generates a presigned URL for uploading a video file
func (s *mediaService) GeneratePresignedUploadURLForVideo(folder, fileName, fileType string) (string, error) {
	return s.s3Client.GeneratePresignedURL(folder, fileName, fileType)
}

// GeneratePresignedUploadURLForImage generates a presigned URL for uploading an image file
func (s *mediaService) GeneratePresignedUploadURLForImage(folder, fileName, fileType string) (string, error) {
	return s.s3Client.GeneratePresignedURL(folder, fileName, fileType)
}

// GeneratePresignedDownloadURLForVideo generates a presigned URL for downloading a video file
func (s *mediaService) GeneratePresignedDownloadURLForVideo(videoID uint64) (string, error) {
	video, err := s.mediaRepo.GetVideoByID(videoID)
	if err != nil {
		return "", err
	}
	if video == nil {
		return "", fmt.Errorf("video not found")
	}

	return s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.FileName, "video/mp4")
}

// GeneratePresignedDownloadURLForImage generates a presigned URL for downloading an image file
func (s *mediaService) GeneratePresignedDownloadURLForImage(videoID uint64) (string, error) {
	video, err := s.mediaRepo.GetVideoByID(videoID)
	if err != nil {
		return "", err
	}
	if video == nil {
		return "", fmt.Errorf("video not found")
	}

	return s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.Image, "image/jpeg")
}
