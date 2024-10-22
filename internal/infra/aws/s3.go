package aws

import (
	"bytes"
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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3ClientInterface interface {
	GeneratePresignedURL(folder string, fileName string, fileType string) (string, error)
	UploadFile(folder string, fileName string, fileType string, fileData []byte) error
}

type S3Client struct {
	Client *s3.Client
	Bucket string
}

func NewS3Client() (S3ClientInterface, error) {
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
func (s *S3Client) GeneratePresignedURL(folder string, fileName string, fileType string) (string, error) {
	log.Info("Folder: ", folder, ", File name: ", fileName)
	if fileName == "" {
		return "", fmt.Errorf("file name must not be empty")
	}

	// Combine folder and fileName to form the S3 key (path to the file)
	fullPath := fileName
	if folder != "" {
		fullPath = folder + "/" + fileName // Add the folder to the file path
	}

	presignClient := s3.NewPresignClient(s.Client)

	reqParams := &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(fullPath), // Use full path (folder + fileName)
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

// UploadFile uploads a file directly to S3
func (s *S3Client) UploadFile(folder string, fileName string, fileType string, fileData []byte) error {
	log.Info("Uploading file to folder: ", folder, ", file name: ", fileName)
	if fileName == "" {
		return fmt.Errorf("file name must not be empty")
	}

	// Combine folder and fileName to form the S3 key (path to the file)
	fullPath := fileName
	if folder != "" {
		fullPath = fmt.Sprintf("%s/%s", folder, fileName)
	}

	// Prepare the PutObject input
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(fullPath),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String(fileType),
		ACL:         types.ObjectCannedACLPrivate, // Set appropriate ACL
	}

	// Perform the upload
	_, err := s.Client.PutObject(context.TODO(), input)
	if err != nil {
		log.Errorf("failed to upload file: %v", err)
		return fmt.Errorf("failed to upload file: %v", err)
	}

	log.Infof("file uploaded successfully: %s", fullPath)
	return nil
}

// DeleteFile deletes a file from the specified S3 folder.
func (s *S3Client) DeleteFile(folder string, fileName string) error {
	if fileName == "" {
		return fmt.Errorf("file name must not be empty")
	}

	// Construct the S3 object key
	fullPath := fileName
	if folder != "" {
		fullPath = fmt.Sprintf("%s/%s", folder, fileName)
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fullPath),
	}

	_, err := s.Client.DeleteObject(context.TODO(), input)
	if err != nil {
		log.Errorf("Failed to delete s3 object: %v", err)
		return fmt.Errorf("Failed to delete s3 object: %v", err)
	}

	// Optionally, wait until the object is deleted
	waiter := s3.NewObjectNotExistsWaiter(s.Client)
	err = waiter.Wait(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(fullPath),
	}, 5*time.Minute)
	if err != nil {
		log.Error("Error waiting for object deletion: ", err)
		return fmt.Errorf("error waiting for object deletion: %v", err)
	}

	log.Info("Successfully deleted file from S3: ", fullPath)
	return nil
}
