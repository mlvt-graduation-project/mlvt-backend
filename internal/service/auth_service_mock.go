package service

import (
	"mlvt/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GenerateToken(user *entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GetUserByToken(tokenStr string) (*entity.User, error) {
	args := m.Called(tokenStr)
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
