package service

import (
	"errors"
	"mlvt/internal/pkg/auth"
	"mlvt/internal/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
	auth     *auth.Auth
}

func NewAuthService(userRepo repository.UserRepository, auth *auth.Auth) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		auth:     auth,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil || user == nil || user.Password != password {
		return "", errors.New("invalid email or password")
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Logout(token string) error {
	return nil
}
