package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	// Define the user input
	input := entity.User{
		FirstName: "Jane",
		LastName:  "Doe",
		UserName:  "janedoe",
		Email:     "jane@example.com",
		Password:  "password123",
	}

	// Mock the RegisterUser method
	mockService.On("RegisterUser", mock.AnythingOfType("*entity.User")).Return(nil)

	// Create the request body
	body, _ := json.Marshal(input)

	// Create a request
	req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a router and register the endpoint
	router := gin.Default()
	router.POST("/users/register", controller.RegisterUser)

	// Perform the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Check the response body
	var resp response.MessageResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestRegisterUser_Failure_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	// Invalid JSON (missing closing brace)
	body := []byte(`{"first_name": "Jane", "email": "jane@example.com"`)

	req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.POST("/users/register", controller.RegisterUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.Error, "unexpected EOF") // Updated expectation

	// Service should not be called
	mockService.AssertNotCalled(t, "RegisterUser", mock.Anything)
}

func TestRegisterUser_Failure_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	input := entity.User{
		FirstName: "Jane",
		LastName:  "Doe",
		UserName:  "janedoe",
		Email:     "jane@example.com",
		Password:  "password123",
	}

	mockService.On("RegisterUser", mock.AnythingOfType("*entity.User")).Return(errors.New("db error"))

	body, _ := json.Marshal(input)

	req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.POST("/users/register", controller.RegisterUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "db error", resp.Error)

	mockService.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "jane@example.com",
		Password: "password123",
	}

	token := "jwt.token.here"

	mockService.On("Login", credentials.Email, credentials.Password).Return(token, nil)

	body, _ := json.Marshal(credentials)

	req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.POST("/users/login", controller.LoginUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp response.TokenResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, token, resp.Token)

	mockService.AssertExpectations(t)
}

func TestLoginUser_Failure_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "jane@example.com",
		Password: "wrongpassword",
	}

	mockService.On("Login", credentials.Email, credentials.Password).Return("", errors.New("invalid credentials"))

	body, _ := json.Marshal(credentials)

	req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.POST("/users/login", controller.LoginUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid credentials", resp.Error)

	mockService.AssertExpectations(t)
}

func TestChangePassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)
	oldPassword := "oldpassword"
	newPassword := "newpassword"

	request := struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	body, _ := json.Marshal(request)

	mockService.On("ChangePassword", userID, oldPassword, newPassword).Return(nil)

	req, err := http.NewRequest(http.MethodPut, "/users/"+strconv.FormatUint(userID, 10)+"/change-password", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.PUT("/users/:user_id/change-password", controller.ChangePassword)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp response.MessageResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Password changed successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestChangePassword_Failure_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	invalidUserID := "abc"
	request := struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}

	body, _ := json.Marshal(request)

	req, err := http.NewRequest(http.MethodPut, "/users/"+invalidUserID+"/change-password", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.PUT("/users/:user_id/change-password", controller.ChangePassword)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid user ID", resp.Error)

	mockService.AssertNotCalled(t, "ChangePassword", mock.Anything, mock.Anything, mock.Anything)
}

func TestChangePassword_Failure_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)
	oldPassword := "oldpassword"
	newPassword := "newpassword"

	request := struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	body, _ := json.Marshal(request)

	mockService.On("ChangePassword", userID, oldPassword, newPassword).Return(errors.New("db error"))

	req, err := http.NewRequest(http.MethodPut, "/users/"+strconv.FormatUint(userID, 10)+"/change-password", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.PUT("/users/:user_id/change-password", controller.ChangePassword)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "db error", resp.Error)

	mockService.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)
	input := entity.User{
		FirstName: "John",
		LastName:  "Doe",
		UserName:  "johndoe",
		Email:     "john@example.com",
		Status:    entity.UserStatusAvailable,
		Premium:   false,
		Role:      "user",
	}

	input.ID = userID

	mockService.On("UpdateUser", &input).Return(nil)

	body, _ := json.Marshal(input)

	req, err := http.NewRequest(http.MethodPut, "/users/"+strconv.FormatUint(userID, 10), bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.PUT("/users/:user_id", controller.UpdateUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp response.MessageResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "User updated successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestUpdateUser_Failure_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	invalidUserID := "abc"
	input := entity.User{
		FirstName: "John",
		LastName:  "Doe",
		UserName:  "johndoe",
		Email:     "john@example.com",
		Status:    entity.UserStatusAvailable,
		Premium:   false,
		Role:      "user",
	}

	body, _ := json.Marshal(input)

	req, err := http.NewRequest(http.MethodPut, "/users/"+invalidUserID, bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.PUT("/users/:user_id", controller.UpdateUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid user ID", resp.Error)

	mockService.AssertNotCalled(t, "UpdateUser", mock.Anything)
}

func TestUpdateAvatar_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize environment configurations
	env.EnvConfig.AvatarFolder = "avatars"

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)
	fileName := "avatar.jpg"
	avatarFolder := "avatars" // Must match env.EnvConfig.AvatarFolder

	// Mock GeneratePresignedAvatarUploadURL
	presignedURL := "https://s3.amazonaws.com/bucket/avatars/avatar.jpg?presigned"
	mockService.On("GeneratePresignedAvatarUploadURL", avatarFolder, fileName, "image/jpeg").Return(presignedURL, nil)

	// Mock UpdateAvatar
	mockService.On("UpdateAvatar", userID, fileName, avatarFolder).Return(nil)

	// Create a request
	req, err := http.NewRequest(http.MethodPut, "/users/1/update-avatar?file_name=avatar.jpg", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.PUT("/users/:user_id/update-avatar", controller.UpdateAvatar)

	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var resp struct {
		AvatarUploadURL string `json:"avatar_upload_url"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, presignedURL, resp.AvatarUploadURL)

	mockService.AssertExpectations(t)
}

func TestUpdateAvatar_Failure_MissingFileName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)

	req, err := http.NewRequest(http.MethodPut, "/users/"+strconv.FormatUint(userID, 10)+"/update-avatar", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.PUT("/users/:user_id/update-avatar", controller.UpdateAvatar)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "file_name is required", resp.Error)

	mockService.AssertNotCalled(t, "GeneratePresignedAvatarUploadURL", mock.Anything, mock.Anything, mock.Anything)
	mockService.AssertNotCalled(t, "UpdateAvatar", mock.Anything, mock.Anything, mock.Anything)
}

func TestLoadAvatar_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)
	expectedURL := "https://s3.amazonaws.com/bucket/avatars/avatar.jpg?presigned"

	mockService.On("GeneratePresignedAvatarDownloadURL", userID).Return(expectedURL, nil)

	req, err := http.NewRequest(http.MethodGet, "/users/"+strconv.FormatUint(userID, 10)+"/avatar", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/users/:user_id/avatar", controller.LoadAvatar)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	assert.Equal(t, expectedURL, rr.Header().Get("Location"))

	mockService.AssertExpectations(t)
}

func TestLoadAvatar_Failure_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	invalidUserID := "abc"

	req, err := http.NewRequest(http.MethodGet, "/users/"+invalidUserID+"/avatar", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/users/:user_id/avatar", controller.LoadAvatar)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid user ID", resp.Error)

	mockService.AssertNotCalled(t, "GeneratePresignedAvatarDownloadURL", mock.Anything)
}

func TestLoadAvatar_Failure_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)

	mockService.On("GeneratePresignedAvatarDownloadURL", userID).Return("", errors.New("s3 error"))

	req, err := http.NewRequest(http.MethodGet, "/users/"+strconv.FormatUint(userID, 10)+"/avatar", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/users/:user_id/avatar", controller.LoadAvatar)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "s3 error", resp.Error)

	mockService.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)
	user := &entity.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		UserName:  "johndoe",
		Email:     "john@example.com",
	}

	mockService.On("GetUserByID", userID).Return(user, nil)

	req, err := http.NewRequest(http.MethodGet, "/users/"+strconv.FormatUint(userID, 10), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/users/:user_id", controller.GetUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp response.UserResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, userID, resp.User.ID)
	assert.Equal(t, "john@example.com", resp.User.Email)

	mockService.AssertExpectations(t)
}

func TestGetUser_Failure_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	invalidUserID := "abc"

	req, err := http.NewRequest(http.MethodGet, "/users/"+invalidUserID, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/users/:user_id", controller.GetUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid user ID", resp.Error)

	mockService.AssertNotCalled(t, "GetUserByID", mock.Anything)
}

func TestGetUser_Failure_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)

	mockService.On("GetUserByID", userID).Return(nil, errors.New("db error"))

	req, err := http.NewRequest(http.MethodGet, "/users/"+strconv.FormatUint(userID, 10), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/users/:user_id", controller.GetUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "db error", resp.Error)

	mockService.AssertExpectations(t)
}

func TestGetAllUsers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	users := []entity.User{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			UserName:  "johndoe",
			Email:     "john@example.com",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Smith",
			UserName:  "janesmith",
			Email:     "jane@example.com",
		},
	}

	mockService.On("GetAllUsers").Return(users, nil)

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/users", controller.GetAllUsers)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp response.UsersResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Users, 2)
	assert.Equal(t, "john@example.com", resp.Users[0].Email)
	assert.Equal(t, "jane@example.com", resp.Users[1].Email)

	mockService.AssertExpectations(t)
}

func TestGetAllUsers_Failure_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	mockService.On("GetAllUsers").Return(nil, errors.New("db error"))

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/users", controller.GetAllUsers)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "internal server error", resp.Error)

	mockService.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)

	mockService.On("DeleteUser", userID).Return(nil)

	req, err := http.NewRequest(http.MethodDelete, "/users/"+strconv.FormatUint(userID, 10), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.DELETE("/users/:user_id", controller.DeleteUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp response.MessageResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "User deleted successfully", resp.Message)

	mockService.AssertExpectations(t)
}

func TestDeleteUser_Failure_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	invalidUserID := "abc"

	req, err := http.NewRequest(http.MethodDelete, "/users/"+invalidUserID, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.DELETE("/users/:user_id", controller.DeleteUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid user ID", resp.Error)

	mockService.AssertNotCalled(t, "DeleteUser", mock.Anything)
}

func TestDeleteUser_Failure_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)

	mockService.On("DeleteUser", userID).Return(errors.New("user not found"))

	req, err := http.NewRequest(http.MethodDelete, "/users/"+strconv.FormatUint(userID, 10), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.DELETE("/users/:user_id", controller.DeleteUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "user not found", resp.Error)

	mockService.AssertExpectations(t)
}

func TestDeleteUser_Failure_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockUserService)
	controller := NewUserController(mockService)

	userID := uint64(1)

	mockService.On("DeleteUser", userID).Return(errors.New("db error"))

	req, err := http.NewRequest(http.MethodDelete, "/users/"+strconv.FormatUint(userID, 10), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := gin.Default()
	router.DELETE("/users/:user_id", controller.DeleteUser)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var resp response.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "internal server error", resp.Error)

	mockService.AssertExpectations(t)
}
