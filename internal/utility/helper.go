package utility

import (
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/pkg/response"
	"sort"

	"github.com/spf13/viper"
)


func GetPostgresConnection () string {
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	sslMode := viper.GetString("SSL_MODE")
	
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslMode)
}


func SortProcessResponseByCreateDate(p []response.ProcessResponse, limit int) []response.ProcessResponse {
	// Sort by created at in descending order
	sort.Slice(p, func(i, j int) bool {
		return p[i].CreatedAt.After(p[j].CreatedAt)
	})

	// Limit the result
	if limit > 0 && limit < len(p) {
		return p[:limit]
	}
	return p
}

func Contains[T comparable](list []T, target T) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func GetMediaTitle(mediaType entity.MediaType, isFullPipeline bool, isOriginalText bool, id int) string {
	var result string
	switch mediaType {
		case entity.MediaTypeAudio:
			result = fmt.Sprintf("Audio - %v", id)
		case entity.MediaTypeVideo:
			result = fmt.Sprintf("Video - %v", id)
		case entity.MediaTypeText:
			if isOriginalText {
				result = fmt.Sprintf("Text - %v", id)
			} else {
				result = fmt.Sprintf("Translated Text - %v", id)
			}
	}
	if isFullPipeline {
		result += " (Full pipeline)"
	}
	return result	
}