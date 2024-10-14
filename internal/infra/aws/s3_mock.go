package aws

import (
	"github.com/stretchr/testify/mock"
)

// MockS3Client is a mock implementation of S3ClientInterface
type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) GeneratePresignedURL(folder string, fileName string, fileType string) (string, error) {
	args := m.Called(folder, fileName, fileType)
	return args.String(0), args.Error(1)
}
