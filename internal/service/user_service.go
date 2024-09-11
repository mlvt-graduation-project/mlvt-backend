package service

import (
	"errors"
	"mlvt/internal/entity"
	"mlvt/internal/infra/reason"
	"mlvt/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

const (
	// ActiveStatus represents the active user status
	ActiveStatus = 1
)

// UserService handles user-related business logic
type UserService struct {
	userRepo    repo.UserRepository
	authService *AuthService
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repo.UserRepository, authService *AuthService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		authService: authService,
	}
}

// RegisterUser handles the registration of a new user
func (s *UserService) RegisterUser(firstName, lastName, email, password string) error {
	// Check if the email is already registered
	existingUser, _ := s.userRepo.GetUserByEmail(email)
	if existingUser != nil {
		return errors.New(reason.EmailAlreadyRegistered.Message())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hashedPassword),
		Status:    ActiveStatus, // Using constant for active status
	}

	return s.userRepo.CreateUser(user)
}

// Login authenticates the user and returns a JWT token
func (s *UserService) Login(email, password string) (string, error) {
	return s.authService.Login(email, password)
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id uint64) (*entity.User, error) {
	return s.userRepo.GetUserByID(id)
}

// UpdateUser updates an existing user's details
func (s *UserService) UpdateUser(id uint64, firstName, lastName, email, password string, status int) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email
	user.Status = status

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	return s.userRepo.UpdateUser(user)
}

// DeleteUser removes a user
func (s *UserService) DeleteUser(id uint64) error {
	return s.userRepo.DeleteUser(id)
}
