package handler

import (
	"mlvt/internal/entity"
	"mlvt/internal/infra/env"
	"mlvt/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

// RegisterUser handles user registration
// @Summary Register a new user
// @Description Creates a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.User true "User data"
// @Success 201 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /users/register [post]
func (h *UserController) RegisterUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// GenerateAvatarDownloadURL generates a presigned URL for downloading the user's avatar
// @Summary Get presigned URL for avatar download
// @Description Generates a presigned URL to download the user's avatar from S3
// @Tags users
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]string "avatar_download_url"
// @Failure 500 {object} map[string]string "error"
// @Router /users/{user_id}/avatar-download-url [get]
func (h *UserController) GenerateAvatarDownloadURL(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	// Call the service to generate the presigned download URL for the avatar
	url, err := h.userService.GeneratePresignedAvatarDownloadURL(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the presigned URL for the avatar image
	c.JSON(http.StatusOK, gin.H{"avatar_download_url": url})
}

// LoginUser handles user login
// @Summary User login
// @Description Logs in a user with their email and password
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body object true "Email and password"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "error"
// @Failure 401 {object} map[string]string "error"
// @Router /users/login [post]
func (h *UserController) LoginUser(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.Login(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ChangePassword handles password changes
// @Summary Change user password
// @Description Allows a user to change their password
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param password body object true "Old and new password"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /users/{user_id}/change-password [put]
func (h *UserController) ChangePassword(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	var request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.ChangePassword(userID, request.OldPassword, request.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// UpdateUser handles user information updates (excluding avatar)
// @Summary Update user information
// @Description Updates the user's information, excluding the avatar
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param user body entity.User true "User data"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /users/{user_id} [put]
func (h *UserController) UpdateUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = userID

	if err := h.userService.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// UpdateAvatar handles avatar updates
// @Summary Update user avatar
// @Description Generates a presigned URL for uploading the user's avatar
// @Tags users
// @Produce json
// @Param user_id path int true "User ID"
// @Param folder query string true "Folder for avatar storage"
// @Param file_name query string true "File name for avatar"
// @Success 200 {object} map[string]string "avatar_upload_url"
// @Failure 500 {object} map[string]string "error"
// @Router /users/{user_id}/update-avatar [put]
func (h *UserController) UpdateAvatar(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	folder := env.EnvConfig.AvatarFolder
	fileName := c.Query("file_name")

	url, err := h.userService.GeneratePresignedAvatarUploadURL(folder, fileName, "image/jpeg")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// You can store the avatar path and folder for the user in the database after a successful upload
	if err := h.userService.UpdateAvatar(userID, fileName, folder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar_upload_url": url})
}

// LoadAvatar loads the user's avatar by redirecting to the presigned URL
// @Summary Load user avatar
// @Description Redirects the client to the presigned URL to download the user's avatar
// @Tags users
// @Produce json
// @Param user_id path int true "User ID"
// @Success 307 {string} string "Redirects to avatar URL"
// @Failure 500 {object} map[string]string "error"
// @Router /users/{user_id}/avatar [get]
func (h *UserController) LoadAvatar(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	// Call the service to generate the presigned download URL for the avatar
	url, err := h.userService.GeneratePresignedAvatarDownloadURL(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Redirect the user to the presigned URL
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Description Fetches a user's details by their ID
// @Tags users
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} entity.User "User data"
// @Failure 500 {object} map[string]string "error"
// @Router /users/{user_id} [get]
func (h *UserController) GetUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetAllUsers retrieves all users
// @Summary Get all users
// @Description Retrieves a list of all users in the system
// @Tags users
// @Produce json
// @Success 200 {array} entity.User "List of users"
// @Failure 500 {object} map[string]string "error"
// @Router /users [get]
func (h *UserController) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// DeleteUser handles user deletion (soft delete)
// @Summary Delete user
// @Description Soft deletes a user by updating their status
// @Tags users
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /users/{user_id} [delete]
func (h *UserController) DeleteUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err := h.userService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
