package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSRepository interface {
	UploadFile(bucketName string, bucketKey string, fileName string) error
	GetFile(fileName string) error
}

type AWSService struct {
	S3Client *s3.S3
}

func NewAWSService(s3Client *s3.S3) *AWSService {
	return &AWSService{
		S3Client: s3Client,
	}
}

func (awsc AWSService) UploadFile(fileName string) error {
	file, err := os.Open(fileName)

	if err != nil {
		log.Println("Error opening file", err)
	} else {
		defer file.Close()

		awsConfig := &aws.Config{
			Region: aws.String("ap-southeast-1"),
		}

		// The session the S3 Uploader will use
		sess := session.Must(session.NewSession(awsConfig))

		// Create an uploader with the session and custom options
		uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
			u.PartSize = 5 * 1024 * 1024 // The minimum/default allowed part size is 5MB
			u.Concurrency = 2            // default is 5
		})

		// Upload the file to S3.
		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("video-bucket-file"),
			Key:    aws.String("example"),
			Body:   file,
		})

		// In case it fails to upload
		if err != nil {
			fmt.Printf("Failed to upload file, %v", err)
			return nil
		}
		fmt.Printf("File uploaded to, %s\n", result.Location)
	}

	log.Println("Successfully uploaded file", fileName)
	return nil
}

func (awsc AWSService) GetFile(item string) error {
	file, err := os.Create(item)
	if err != nil {
		log.Println("Error creating file", err)
	}

	defer file.Close()

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	})

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String("bucket-video-file"),
			Key:    aws.String(item),
		})

	if err != nil {
		log.Println("Error downloading file", err)
	}

	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")
	return nil
}
