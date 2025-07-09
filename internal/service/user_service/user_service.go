package user_service

import (
	"context"
	"errors"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/repo/user_repo"
	"mlvt/internal/service/auth_service"
	"mlvt/internal/service/email_service"
	"mlvt/internal/service/traffic_service"
	"mlvt/internal/utility"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(user *entity.User) error
	Login(email, password string) (string, uint64, entity.UserPermission, error)
	ChangePassword(userID uint64, oldPassword, newPassword string) error
	UpdateUser(user *entity.User) error
	UpdateAvatar(userID uint64, avatarPath, avatarFolder string) error
	VerifyAccountSignUp(username string, token string) error
	ResendValidationEmail(username string, email string) error
	GetUserByID(userID uint64) (*entity.User, error)
	GetAllUsers() ([]entity.User, error)
	DeleteUser(userID uint64) error
	GeneratePresignedAvatarUploadURL(folder, fileName, fileType string) (string, error)
	GeneratePresignedAvatarDownloadURL(userID uint64) (string, error)
}

type userService struct {
	repo           user_repo.UserRepository
	s3Client       aws.S3ClientInterface
	auth           auth_service.AuthServiceInterface
	trafficService traffic_service.TrafficService
	emailService   email_service.EmailService
}

func NewUserService(
	repo user_repo.UserRepository,
	s3Client aws.S3ClientInterface,
	auth auth_service.AuthServiceInterface,
	trafficService traffic_service.TrafficService,
	emailService email_service.EmailService,
) UserService {
	return &userService{
		repo:           repo,
		s3Client:       s3Client,
		auth:           auth,
		trafficService: trafficService,
		emailService:   emailService,
	}
}

// RegisterUser creates a new user with hashed password
func (s *userService) RegisterUser(user *entity.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.Status = entity.UserStatusActive
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	ctx := context.Background()
	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.CreateAccountAction,
		Description: "new user created",
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: create new user account")
	}

	// set user status to pending
	user.Status = entity.UserStatusPending
	if user.Role == "" {
		user.Role = entity.UserRole
	}

	checkUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		log.Errorf("error checking existing email when registing user. Error: %v", err)
		return fmt.Errorf("failed to register user")
	}
	if checkUser != nil {
		return fmt.Errorf("email has already been register")
	}

	// Create user to database
	err = s.repo.CreateUser(user)
	if err != nil {
		log.Errorf("failed to create user to database")
		return err
	}

	expiredTime := utility.SetExpireTime(180)

	// Encrypt token for account sign up
	token, err := utility.EncryptToken(user.UserName, expiredTime)
	if err != nil {
		log.Errorf("failed to encrypt token", err)
		return err
	}

	// Create html body for account sign up
	htmlBody, err := s.emailService.CreateAccountSignUpEmail(user.UserName, token)
	if err != nil {
		log.Errorf("failed to created html body for account sign up", err)
		return err
	}

	// Send email to user
	err = s.emailService.SendHTMLEmail("MLVT Account created successfully", htmlBody, user.Email)
	if err != nil {
		log.Errorf("failed to send html email to account", err)
	}

	return nil
}

// VerifyAccountSignUp verifies the account sign up
func (s *userService) VerifyAccountSignUp(email string, token string) error {
	user := &entity.User{
		Email:  email,
		Status: entity.UserStatusPending,
	}

	// Get user by email
	user, err := s.repo.GetUserByCondition(user)
	if err != nil {
		log.Errorf("Invalid user credential", err)
		return err
	}

	// Decrypt token
	username, expireDate, err := utility.DecryptToken(token)
	if err != nil {
		log.Errorf("failed to decrypt token", err)
		return err
	} else if username != user.UserName {
		log.Errorf("Invalid user credential", err)
	}

	// Check if the token is expired
	if time.Now().After(expireDate) {
		return errors.New("token expired")
	}

	user.Status = entity.UserStatusActive
	err = s.repo.UpdateUser(user)
	if err != nil {
		return errors.New("failed to update user status")
	}

	return nil
}

func (s *userService) ResendValidationEmail(username string, email string) error {
	user := &entity.User{}
	if username != "" {
		user.UserName = username
	}
	if email != "" {
		user.Email = email
	}
	// Get user by username
	user, err := s.repo.GetUserByCondition(user)
	if err != nil {
		log.Errorf("Invalid user credential", err)
		return err
	}

	expiredTime := utility.SetExpireTime(180)

	// Encrypt token for account sign up
	token, err := utility.EncryptToken(user.UserName, expiredTime)
	if err != nil {
		log.Errorf("failed to encrypt token", err)
		return err
	}

	// Create html body for account sign up
	htmlBody, err := s.emailService.CreateAccountSignUpEmail(user.UserName, token)
	if err != nil {
		log.Errorf("failed to created html body for account sign up", err)
		return err
	}

	// Send email to user
	err = s.emailService.SendHTMLEmail("MLVT Account created successfully", htmlBody, user.Email)
	if err != nil {
		log.Errorf("failed to send html email to account", err)
		return err
	}

	return nil

}

// Login handles user login
func (s *userService) Login(email, password string) (string, uint64, entity.UserPermission, error) {
	ctx := context.Background()
	token, userID, role, err := s.auth.Login(email, password)

	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.LoginAction,
		Description: "user login",
		UserID:      userID,
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: login account")
	}

	return token, userID, role, err
}

// ChangePassword changes a user's password
func (s *userService) ChangePassword(userID uint64, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Compare old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("old password does not match")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.ChangePasswordAction,
		Description: "user change password",
		UserID:      userID,
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: change password account")
	}

	return s.repo.UpdateUserPassword(userID, string(hashedPassword))
}

// UpdateUser updates user information (except avatar)
func (s *userService) UpdateUser(user *entity.User) error {
	user.UpdatedAt = time.Now()

	ctx := context.Background()
	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.UpdateProfileAction,
		Description: "user update profile",
		UserID:      user.ID,
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: update profile")
	}
	return s.repo.UpdateUser(user)
}

// UpdateAvatar updates the user's avatar
func (s *userService) UpdateAvatar(userID uint64, avatarPath, avatarFolder string) error {
	ctx := context.Background()
	if _, err := s.trafficService.CreateTraffic(ctx, entity.Traffic{
		ActionType:  entity.UploadAvatarAction,
		Description: "user update avatar",
		UserID:      userID,
		Timestamp:   time.Now().Unix(),
	}); err != nil {
		log.Errorf("failed to log traffic: update avatar")
	}
	return s.repo.UpdateUserAvatar(userID, avatarPath, avatarFolder)
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(userID uint64) (*entity.User, error) {
	return s.repo.GetUserByID(userID)
}

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers() ([]entity.User, error) {
	return s.repo.GetAllUsers()
}

// DeleteUser soft deletes a user by setting their status to "deleted"
func (s *userService) DeleteUser(userID uint64) error {
	return s.repo.DeleteUser(userID)
}

// GeneratePresignedAvatarUploadURL generates a presigned URL for uploading an avatar
func (s *userService) GeneratePresignedAvatarUploadURL(folder, fileName, fileType string) (string, error) {
	return s.s3Client.GeneratePresignedURL(folder, fileName, fileType)
}

// GeneratePresignedAvatarDownloadURL generates a presigned URL for downloading the user's avatar
func (s *userService) GeneratePresignedAvatarDownloadURL(userID uint64) (string, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	if user.Avatar == "" || user.AvatarFolder == "" {
		return "", errors.New("avatar not found for this user")
	}

	// Generate the presigned URL for the avatar image
	url, err := s.s3Client.GeneratePresignedDownloadURL(user.AvatarFolder, user.Avatar, "image/jpeg")
	if err != nil {
		return "", err
	}

	return url, nil
}
