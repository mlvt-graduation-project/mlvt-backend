package service

import (
	"errors"
	"testing"
	"time"

	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	user := &entity.User{
		FirstName: "John",
		LastName:  "Doe",
		UserName:  "johndoe",
		Email:     "john@example.com",
		Password:  "password123", // Plain password
	}

	// Expect CreateUser to be called with the user (password should be hashed)
	mockRepo.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(nil)

	err := userService.RegisterUser(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Password) // Password should be hashed

	// Verify password is hashed
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_Failure_HashPassword(t *testing.T) {
	// Note: bcrypt.GenerateFromPassword is not easily mockable without additional interfaces.
	// As an alternative, you can skip this test or refactor your code to allow mocking bcrypt.

	// For demonstration, we'll simulate a failure in CreateUser
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	user := &entity.User{
		FirstName: "John",
		LastName:  "Doe",
		UserName:  "johndoe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(errors.New("db error"))

	err := userService.RegisterUser(user)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	email := "john@example.com"
	password := "password123"
	token := "jwt.token.here"

	mockAuth.On("Login", email, password).Return(token, nil)

	returnedToken, err := userService.Login(email, password)
	assert.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	mockAuth.AssertExpectations(t)
}

func TestLogin_Failure_InvalidCredentials(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	email := "john@example.com"
	password := "wrongpassword"

	mockAuth.On("Login", email, password).Return("", errors.New("invalid credentials"))

	returnedToken, err := userService.Login(email, password)
	assert.Error(t, err)
	assert.Equal(t, "", returnedToken)
	assert.Equal(t, "invalid credentials", err.Error())

	mockAuth.AssertExpectations(t)
}

func TestChangePassword_Success(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)
	oldPassword := "oldpassword"
	newPassword := "newpassword"

	// Hash the old password
	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

	user := &entity.User{
		ID:       userID,
		Password: string(hashedOldPassword),
	}

	mockRepo.On("GetUserByID", userID).Return(user, nil)
	mockRepo.On("UpdateUserPassword", userID, mock.AnythingOfType("string")).Return(nil)

	err := userService.ChangePassword(userID, oldPassword, newPassword)
	assert.NoError(t, err)

	// Verify the new password is hashed and updated
	mockRepo.AssertCalled(t, "UpdateUserPassword", userID, mock.MatchedBy(func(hashed string) bool {
		err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(newPassword))
		return err == nil
	}))
	mockRepo.AssertExpectations(t)
}

func TestChangePassword_Failure_WrongOldPassword(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)
	oldPassword := "correctpassword"
	wrongOldPassword := "wrongpassword"
	newPassword := "newpassword"

	// Hash the correct old password
	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

	user := &entity.User{
		ID:       userID,
		Password: string(hashedOldPassword),
	}

	mockRepo.On("GetUserByID", userID).Return(user, nil)

	err := userService.ChangePassword(userID, wrongOldPassword, newPassword)
	assert.Error(t, err)
	assert.Equal(t, "old password does not match", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestChangePassword_Failure_UserNotFound(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)
	oldPassword := "oldpassword"
	newPassword := "newpassword"

	mockRepo.On("GetUserByID", userID).Return(nil, errors.New("user not found"))

	err := userService.ChangePassword(userID, oldPassword, newPassword)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	user := &entity.User{
		ID:        1,
		FirstName: "Jane",
		LastName:  "Doe",
		UserName:  "janedoe",
		Email:     "jane@example.com",
		Status:    entity.UserStatusAvailable,
		Premium:   true,
		Role:      "admin",
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UpdateUser", user).Return(nil)

	err := userService.UpdateUser(user)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Failure(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	user := &entity.User{
		ID:        1,
		FirstName: "Jane",
		LastName:  "Doe",
		UserName:  "janedoe",
		Email:     "jane@example.com",
		Status:    entity.UserStatusAvailable,
		Premium:   true,
		Role:      "admin",
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UpdateUser", user).Return(errors.New("update error"))

	err := userService.UpdateUser(user)
	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestUpdateAvatar_Success(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)
	avatarPath := "avatar_new.jpg"
	avatarFolder := "avatars_new"

	mockRepo.On("UpdateUserAvatar", userID, avatarPath, avatarFolder).Return(nil)

	err := userService.UpdateAvatar(userID, avatarPath, avatarFolder)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)
	user := &entity.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
	}

	mockRepo.On("GetUserByID", userID).Return(user, nil)

	returnedUser, err := userService.GetUserByID(userID)
	assert.NoError(t, err)
	assert.Equal(t, user, returnedUser)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Failure(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)

	mockRepo.On("GetUserByID", userID).Return(nil, errors.New("user not found"))

	returnedUser, err := userService.GetUserByID(userID)
	assert.Error(t, err)
	assert.Nil(t, returnedUser)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_Success(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	users := []entity.User{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Smith",
		},
	}

	mockRepo.On("GetAllUsers").Return(users, nil)

	returnedUsers, err := userService.GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, users, returnedUsers)

	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_Failure(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	mockRepo.On("GetAllUsers").Return(nil, errors.New("db error"))

	returnedUsers, err := userService.GetAllUsers()
	assert.Error(t, err)
	assert.Nil(t, returnedUsers)
	assert.Equal(t, "db error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)

	mockRepo.On("DeleteUser", userID).Return(nil)

	err := userService.DeleteUser(userID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Failure(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)

	mockRepo.On("DeleteUser", userID).Return(errors.New("user not found"))

	err := userService.DeleteUser(userID)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGeneratePresignedAvatarUploadURL_Success(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	folder := "avatars"
	fileName := "avatar_new.jpg"
	fileType := "image/jpeg"
	expectedURL := "https://s3.amazonaws.com/bucket/avatars/avatar_new.jpg?presigned"

	mockS3.On("GeneratePresignedURL", folder, fileName, fileType).Return(expectedURL, nil)

	url, err := userService.GeneratePresignedAvatarUploadURL(folder, fileName, fileType)
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)

	mockS3.AssertExpectations(t)
}

func TestGeneratePresignedAvatarDownloadURL_Success(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)
	user := &entity.User{
		ID:           userID,
		Avatar:       "avatar.jpg",
		AvatarFolder: "avatars",
	}

	expectedURL := "https://s3.amazonaws.com/bucket/avatars/avatar.jpg?presigned"

	mockRepo.On("GetUserByID", userID).Return(user, nil)
	mockS3.On("GeneratePresignedURL", user.AvatarFolder, user.Avatar, "image/jpeg").Return(expectedURL, nil)

	url, err := userService.GeneratePresignedAvatarDownloadURL(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)

	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
}

func TestGeneratePresignedAvatarDownloadURL_Failure_UserNotFound(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)

	mockRepo.On("GetUserByID", userID).Return(nil, errors.New("user not found"))

	url, err := userService.GeneratePresignedAvatarDownloadURL(userID)
	assert.Error(t, err)
	assert.Equal(t, "", url)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGeneratePresignedAvatarDownloadURL_Failure_AvatarNotFound(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)
	user := &entity.User{
		ID:           userID,
		Avatar:       "",
		AvatarFolder: "",
	}

	mockRepo.On("GetUserByID", userID).Return(user, nil)

	url, err := userService.GeneratePresignedAvatarDownloadURL(userID)
	assert.Error(t, err)
	assert.Equal(t, "", url)
	assert.Equal(t, "avatar not found for this user", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGeneratePresignedAvatarDownloadURL_Failure_S3Error(t *testing.T) {
	mockRepo := new(repo.MockUserRepository)
	mockS3 := new(aws.MockS3Client)
	mockAuth := new(MockAuthService)

	userService := NewUserService(mockRepo, mockS3, mockAuth)

	userID := uint64(1)
	user := &entity.User{
		ID:           userID,
		Avatar:       "avatar.jpg",
		AvatarFolder: "avatars",
	}

	mockRepo.On("GetUserByID", userID).Return(user, nil)
	mockS3.On("GeneratePresignedURL", user.AvatarFolder, user.Avatar, "image/jpeg").Return("", errors.New("s3 error"))

	url, err := userService.GeneratePresignedAvatarDownloadURL(userID)
	assert.Error(t, err)
	assert.Equal(t, "", url)
	assert.Equal(t, "s3 error", err.Error())

	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
}
