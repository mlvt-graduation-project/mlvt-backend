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

		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String("ap-southeast-1"),
		})

		uploader := s3manager.NewUploader(sess)

		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("bucket-video-file"),
			Key:    aws.String("example"),
			Body:   file,
		})

		if err != nil {
			log.Println("Error uploading file", err)
		}
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
