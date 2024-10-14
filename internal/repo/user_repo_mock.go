package repo

import (
	"mlvt/internal/entity"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository mocks the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByID(userID uint64) (*entity.User, error) {
	args := m.Called(userID)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetAllUsers() ([]entity.User, error) {
	args := m.Called()
	if users, ok := args.Get(0).([]entity.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUserPassword(userID uint64, hashedPassword string) error {
	args := m.Called(userID, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserAvatar(userID uint64, avatarPath, avatarFolder string) error {
	args := m.Called(userID, avatarPath, avatarFolder)
	return args.Error(0)
}
