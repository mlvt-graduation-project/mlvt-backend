package media_service

import (
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/pkg/response"
	"mlvt/internal/utility"
)

func (s *mediaService) GetAllMedia(
	userID uint64,
	searchKey string,
	limit int,
	offset int,
	status []entity.StatusEntity,
	mediaType []entity.MediaType,
) (response.ProcessResponse, error) {
	var res response.ProcessResponse
	videos, audios, transcriptions, totalCount, err := s.mediaRepo.GetAllMedia(userID, searchKey, limit, offset, status, mediaType)
	if err != nil {
		return res, err
	}

	var videoWithURLsList []response.ListVideosByUserIDResponse
	for _, video := range videos {
		// Generate the presigned URL for the video
		videoURL := ""
		if video.Status != entity.StatusFailed && video.Status != entity.StatusProcessing {
			videoURL, err = s.s3Client.GeneratePresignedDownloadURL(video.Folder, video.FileName, "video/mp4")
			if err != nil {
				return res, fmt.Errorf("failed to generate presigned video URL for video ID %d: %v", video.ID, err)
			}
		}

		// Generate the presigned URL for the image/frame
		imageURL, err := s.s3Client.GeneratePresignedDownloadURL(env.EnvConfig.VideoFramesFolder, video.Image, "image/jpeg")
		if err != nil {
			return res, fmt.Errorf("failed to generate presigned image URL for video ID %d: %v", video.ID, err)
		}

		// Append to the list
		videoWithURLsList = append(videoWithURLsList, response.ListVideosByUserIDResponse{
			Video:    video,
			VideoURL: videoURL,
			ImageURL: imageURL,
		})
	}
	videoList := utility.VideoResponseListToProcessResponseList(videoWithURLsList)
	res.ProcessList = append(res.ProcessList, videoList...)
	audioList := utility.AudioListToProcessResponseList(audios)
	res.ProcessList = append(res.ProcessList, audioList...)
	transcriptionList := utility.TranscriptionListToProcessResponseList(transcriptions)
	res.ProcessList = append(res.ProcessList, transcriptionList...)
	res.TotalCount = totalCount
	res.ProcessList = utility.SortProcessResponseByCreateDate(res.ProcessList, limit)
	return res, nil
}
