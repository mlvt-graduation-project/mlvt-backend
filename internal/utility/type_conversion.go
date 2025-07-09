package utility

import (
	"mlvt/internal/entity"
	"mlvt/internal/pkg/response"
	"strconv"
)

func AudioToProcessResponse(a entity.Audio) *entity.Process {
	return &entity.Process{
		ID:                        strconv.FormatUint(a.ID, 10),
		UserID:                    a.UserID,
		MediaType:                 "audio",
		TranslatedTranscriptionID: a.TranscriptionID,
		AudioID:                   a.ID,
		Status:                    a.Status,
		CreatedAt:                 a.CreatedAt,
		UpdatedAt:                 a.UpdatedAt,
		Language:                  a.Lang,
		Title:                     a.Title,
	}
}

func TranscriptionToProcessResponse(transcription entity.Transcription) *entity.Process {
	return &entity.Process{
		ID:                      strconv.FormatUint(transcription.ID, 10),
		UserID:                  transcription.UserID,
		MediaType:               "transcription",
		OriginalVideoID:         transcription.VideoID,
		OriginalTranscriptionID: transcription.OriginalTranscriptionID,
		Status:                  transcription.Status,
		CreatedAt:               transcription.CreatedAt,
		UpdatedAt:               transcription.UpdatedAt,
		Language:                transcription.Lang,
		Title:                   transcription.Title,
	}
}

func ProgressResponseListToProcessResponseList(p []response.ProgressResponse) []entity.Process {
	var res []entity.Process
	for _, ele := range p {
		res = append(res, *ele.ToProcessResponse())
	}
	return res
}

func TranscriptionListToProcessResponseList(p []entity.Transcription) []entity.Process {
	var res []entity.Process
	for _, ele := range p {
		res = append(res, *TranscriptionToProcessResponse(ele))
	}
	return res
}

func AudioListToProcessResponseList(p []entity.Audio) []entity.Process {
	var res []entity.Process
	for _, ele := range p {
		res = append(res, *AudioToProcessResponse(ele))
	}
	return res
}

func VideoResponseListToProcessResponseList(p []response.ListVideosByUserIDResponse) []entity.Process {
	var res []entity.Process
	for _, ele := range p {
		res = append(res, *ele.ToProcessResponse())
	}
	return res
}
