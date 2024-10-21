package initialize

import (
	"fmt"
	"mlvt/internal/infra/aws"
	"mlvt/internal/infra/zap-logging/log"
)

// InitAWS initializes AWS services like S3.
func InitAWS() (aws.S3ClientInterface, error) {
	s3Client, err := aws.NewS3Client()
	if err != nil {
		log.Errorf("Failed to initialize AWS S3 client: %v", err)
		return nil, fmt.Errorf("failed to initialize AWS S3 client: %w", err)
	}
	return s3Client, nil
}
