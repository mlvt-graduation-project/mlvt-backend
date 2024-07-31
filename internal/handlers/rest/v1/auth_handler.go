package handlers

import (
	"net/http"

	utils "mlvt/internal/pkg"
	"mlvt/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := utils.ReadJSON(c, &req); err != nil {
		utils.ErrorJSON(c, err, http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrorJSON(c, err, http.StatusUnauthorized)
		return
	}

	utils.WriteJSON(c, http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	utils.WriteJSON(c, http.StatusOK, gin.H{"message": "successfully logged out"})
}
