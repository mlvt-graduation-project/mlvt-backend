package request

import "mlvt/internal/entity"



type ProcessRequest struct {
	SearchKey string 					`json:"search_key"`
	ProjectType []entity.ProgressType	`json:"project_type"`
	MediaType []entity.MediaType		`json:"media_type"`
	Status []entity.StatusEntity		`json:"status"`
	Offset int							`json:"offset"`
	Limit int							`json:"limit"`
}