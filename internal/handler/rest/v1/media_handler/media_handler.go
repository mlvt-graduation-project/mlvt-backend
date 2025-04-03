package media_handler

import (
	"mlvt/internal/service/media_service"
)

type MediaController struct {
	mediaService media_service.MediaService
}

func NewMediaController(
	mediaService media_service.MediaService,
) *MediaController {
	return &MediaController{
		mediaService: mediaService,
	}
}
