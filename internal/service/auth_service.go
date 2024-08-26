package service

import (
	"errors"
	"mlvt/internal/entity"
	"mlvt/internal/repo"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles user authentication
type AuthService struct {
	userRepo  repo.UserRepository
	secretKey string
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo repo.UserRepository, secretKey string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		secretKey: secretKey,
	}
}

// Login authenticates the user and returns a JWT token
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

// GenerateToken creates a JWT token for a user
func (s *AuthService) GenerateToken(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"email":  user.Email,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // Token valid for 72 hours
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserByToken extracts user information from a JWT token
func (s *AuthService) GetUserByToken(tokenStr string) (*entity.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Safely assert types from claims
	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		return nil, errors.New("invalid userID type in token")
	}
	userID := uint64(userIDFloat)

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
