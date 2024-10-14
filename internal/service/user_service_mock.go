package service

import (
	"mlvt/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) ChangePassword(userID uint64, oldPassword, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserService) UpdateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) UpdateAvatar(userID uint64, avatarPath, avatarFolder string) error {
	args := m.Called(userID, avatarPath, avatarFolder)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(userID uint64) (*entity.User, error) {
	args := m.Called(userID)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) GetAllUsers() ([]entity.User, error) {
	args := m.Called()
	if users, ok := args.Get(0).([]entity.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) DeleteUser(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserService) GeneratePresignedAvatarUploadURL(folder, fileName, fileType string) (string, error) {
	args := m.Called(folder, fileName, fileType)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) GeneratePresignedAvatarDownloadURL(userID uint64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}
