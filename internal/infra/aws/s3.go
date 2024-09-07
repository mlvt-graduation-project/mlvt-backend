package aws

import (
	"context"
	"fmt"
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/reason"
	"mlvt/internal/infra/zap-logging/log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client *s3.Client
	Bucket string
}

func NewS3Client() (*S3Client, error) {
	// Load the default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(env.EnvConfig.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			env.EnvConfig.AWSAccessKeyID,
			env.EnvConfig.AWSSecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf(reason.UnableToLoadAWSConfig.Message()+": %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)
	bucket := env.EnvConfig.AWSBucket
	log.Info("Using bucket: ", bucket)

	return &S3Client{Client: client, Bucket: bucket}, nil
}

// GeneratePresignedURL generates a presigned URL for uploading a file to S3
func (s *S3Client) GeneratePresignedURL(fileName string, fileType string) (string, error) {
	log.Info("File name: ", fileName)
	if fileName == "" {
		return "", fmt.Errorf("file name must not be empty")
	}

	presignClient := s3.NewPresignClient(s.Client)

	reqParams := &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(fileName),
		ContentType: aws.String(fileType),
	}

	// Use functional options to set the expiration time
	presignReq, err := presignClient.PresignPutObject(context.TODO(), reqParams, func(o *s3.PresignOptions) {
		o.Expires = 15 * time.Minute // Set the expiration time for the presigned URL
	})
	if err != nil {
		log.Error(reason.FailedToPresignPutObjectRequest.Message()+": ", err)
		return "", fmt.Errorf(reason.FailedToPresignPutObjectRequest.Message()+", %v", err)
	}

	log.Info(reason.GeneratedPresignedURL.Message()+": ", presignReq.URL)
	return presignReq.URL, nil
}
