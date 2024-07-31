package service

import (
	"errors"
	"mlvt/internal/models"
	"mlvt/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) RegisterUser(firstName, lastName, email, password string) (*models.User, error) {
	//check if the email is exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	user := &models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
	err = s.userRepo.Create(user)
	return user, err
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *UserService) GetUserByID(id uint64) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(user *models.User) error {
	return s.userRepo.Delete(user)
}
