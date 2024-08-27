package handler

import (
	"errors"
	"mlvt/internal/pkg/json"
	"mlvt/internal/schema"
	"mlvt/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController manages user-related requests
type UserController struct {
	userService *service.UserService
}

// NewUserController creates a new UserController
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// RegisterUser handles user registration
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body schema.RegisterUserRequest true "User registration details"
// @Success 201 {object} map[string]string "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/register [post]
func (uc *UserController) RegisterUser(ctx *gin.Context) {
	var req schema.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		json.ErrorJSON(ctx, err, http.StatusBadRequest)
		return
	}

	err := uc.userService.RegisterUser(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		json.ErrorJSON(ctx, err, http.StatusInternalServerError)
		return
	}

	json.WriteJSON(ctx, http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles user login
// @Summary User login
// @Description Logs in a user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param login body schema.LoginUserRequest true "Login credentials"
// @Success 200 {object} schema.LoginResponse "Token"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /users/login [post]
func (uc *UserController) Login(ctx *gin.Context) {
	var req schema.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		json.ErrorJSON(ctx, err, http.StatusBadRequest)
		return
	}

	token, err := uc.userService.Login(req.Email, req.Password)
	if err != nil {
		json.ErrorJSON(ctx, err, http.StatusUnauthorized)
		return
	}

	json.WriteJSON(ctx, http.StatusOK, schema.LoginResponse{Token: token})
}

// GetUser fetches user details by ID
// @Summary Get user by ID
// @Description Get details of a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entity.User "User details"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 404 {object} map[string]string "User not found"
// @Router /users/{id} [get]
func (uc *UserController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		json.ErrorJSON(ctx, errors.New("Invalid user ID"), http.StatusBadRequest)
		return
	}

	user, err := uc.userService.GetUser(id)
	if err != nil {
		json.ErrorJSON(ctx, errors.New("User not found"), http.StatusNotFound)
		return
	}

	json.WriteJSON(ctx, http.StatusOK, user)
}

// UpdateUser handles updating user details
// @Summary Update user details
// @Description Update the details of an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body schema.RegisterUserRequest true "Updated user details"
// @Success 200 {object} map[string]string "User updated successfully"
// @Failure 400 {object} map[string]string "Invalid user ID or request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [put]
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		json.ErrorJSON(ctx, errors.New("Invalid user ID"), http.StatusBadRequest)
		return
	}

	var req schema.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		json.ErrorJSON(ctx, err, http.StatusBadRequest)
		return
	}

	err = uc.userService.UpdateUser(id, req.FirstName, req.LastName, req.Email, req.Password, 1) // Assuming 1 is active status
	if err != nil {
		json.ErrorJSON(ctx, err, http.StatusInternalServerError)
		return
	}

	json.WriteJSON(ctx, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

// DeleteUser handles deleting a user by ID
// @Summary Delete user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		json.ErrorJSON(ctx, errors.New("Invalid user ID"), http.StatusBadRequest)
		return
	}

	err = uc.userService.DeleteUser(id)
	if err != nil {
		json.ErrorJSON(ctx, err, http.StatusInternalServerError)
		return
	}

	json.WriteJSON(ctx, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
