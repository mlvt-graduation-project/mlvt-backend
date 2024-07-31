package handlers

import (
	"mlvt/internal/models"
	utils "mlvt/internal/pkg"
	"mlvt/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	if err := utils.ReadJSON(c, &req); err != nil {
		utils.ErrorJSON(c, err, http.StatusBadRequest)
		return
	}

	user, err := h.userService.RegisterUser(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(c, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusNotFound)
		return
	}

	utils.WriteJSON(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req models.User
	if err := utils.ReadJSON(c, &req); err != nil {
		utils.ErrorJSON(c, err, http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateUser(&req); err != nil {
		utils.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(c, http.StatusOK, req)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusNotFound)
		return
	}

	if err := h.userService.DeleteUser(user); err != nil {
		utils.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(c, http.StatusOK, gin.H{"message": "user deleted"})
}
