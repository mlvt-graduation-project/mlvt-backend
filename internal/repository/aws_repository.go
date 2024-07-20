package repository

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSRepository interface {
	UploadFile(bucketName string, bucketKey string, fileName string) error
	GetFile(fileName string) error
}

type AWSService struct {
	S3Client *s3.Client
}

func NewAWSService(s3Client *s3.Client) *AWSService {
	return &AWSService{
		S3Client: s3Client,
	}
}

func (awsc AWSService) UploadFile(bucketName string, bucketKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Error opening file", err)
	} else {
		defer file.Close()

		_, err := awsc.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(bucketKey),
			Body:   file,
		})

		if err != nil {
			log.Println("Error uploading file", err)
		}
	}
	log.Println("Successfully uploaded file", fileName)
	return nil
}

func (awsc AWSService) GetFile(fileName string) error {
	// err := awsc.GetFile(fileName)
	// if err != nil {
	// 	log.Println("Error getting file", err)
	// } else {
	// 	log.Println("Successfully got file", fileName)
	// }
	// return err
	return nil
}
