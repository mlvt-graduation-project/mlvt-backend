package token_claim_handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"mlvt/internal/pkg/response"
	"mlvt/internal/repo/token_claim_repo"
	"mlvt/internal/service/token_claim_service"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
	svc token_claim_service.TokenService
}

func New(s token_claim_service.TokenService) *TokenController { return &TokenController{svc: s} }

// parseUID converts a user_id query/path string to uint64.
func parseUID(idStr string) (uint64, error) {
	return strconv.ParseUint(idStr, 10, 64)
}

/* ---------- CLAIM ENDPOINTS ---------- */

// ClaimDaily godoc
// @Summary Claim daily free tokens
// @Description Adds 5 tokens to the user's wallet if not yet claimed today
// @Tags tokens
// @Accept json
// @Produce json
// @Param user_id query uint64 true "User ID"
// @Success 200 {object} response.MessageResponse "message"
// @Failure 400 {object} response.ErrorResponse "invalid or already claimed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /token/daily [post]
func (h *TokenController) ClaimDaily(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := parseUID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user_id"})
		return
	}

	if err := h.svc.ClaimDaily(context.Background(), userID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.MessageResponse{Message: "+5 tokens credited"})
}

// ClaimPremium godoc
// @Summary Claim daily premium tokens
// @Description Adds 20 tokens to the user's wallet if premium and not yet claimed today
// @Tags tokens
// @Accept json
// @Produce json
// @Param user_id query uint64 true "User ID"
// @Success 200 {object} response.MessageResponse "message"
// @Failure 400 {object} response.ErrorResponse "not premium or already claimed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /token/premium [post]
func (h *TokenController) ClaimPremium(c *gin.Context) {
	uid, err := parseUID(c.Query("user_id"))
	if err != nil {
		c.JSON(400, response.ErrorResponse{Error: "invalid user_id"})
		return
	}

	days, err := h.svc.ClaimPremium(c.Request.Context(), uid)
	switch err {
	case nil:
		total := days * 20
		msg := fmt.Sprintf("+%d premium tokens credited for %d days", total, days)
		c.JSON(200, response.MessageResponse{Message: msg})
	case token_claim_repo.ErrNotPremium:
		c.JSON(400, response.ErrorResponse{Error: "not premium or expired"})
	case token_claim_repo.ErrAlreadyClaimed:
		c.JSON(400, response.ErrorResponse{Error: "already claimed today"})
	default:
		c.JSON(500, response.ErrorResponse{Error: err.Error()})
	}
}

/* ---------- LIST / CHECK ENDPOINTS ---------- */

// ListClaims godoc
// @Summary List all daily token claim records
// @Description Retrieves all claim logs, ordered by newest first
// @Tags tokens
// @Produce json
// @Success 200 {object} response.DailyTokenClaimsResponse "claims"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /token/claims [get]
func (h *TokenController) ListClaims(c *gin.Context) {
	claims, err := h.svc.ListClaims(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.DailyTokenClaimsResponse{Claims: claims})
}

// AddPremium godoc
// @Summary   Grant 30-day premium to a user
// @Tags      tokens
// @Accept    json
// @Produce   json
// @Param     user_id  query  uint64  true  "User ID"
// @Success   200      {object}  response.MessageResponse  "message"
// @Failure   400      {object}  response.ErrorResponse   "invalid or missing user_id"
// @Failure   500      {object}  response.ErrorResponse   "server error"
// @Router    /premium/add [post]
func (h *TokenController) AddPremium(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := parseUID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user_id"})
		return
	}

	if err := h.svc.AddPremium(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.MessageResponse{Message: "premium granted for 30 days"})
}

// ListPremium godoc
// @Summary List all active premium users
// @Description Returns users with premium expiry time in the future
// @Tags tokens
// @Produce json
// @Success 200 {object} response.PremiumUsersResponse "users"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /premium/list [get]
func (h *TokenController) ListPremium(c *gin.Context) {
	users, err := h.svc.ListPremiumUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.PremiumUsersResponse{Users: users})
}

// CheckPremium godoc
// @Summary Check premium status of a user
// @Description Returns whether the user is premium and their expiry time
// @Tags tokens
// @Produce json
// @Param user_id path uint64 true "User ID"
// @Success 200 {object} map[string]interface{} "example: {\"premium\": true}"
// @Failure 400 {object} response.ErrorResponse "invalid user_id"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /premium/{user_id} [get]
func (h *TokenController) CheckPremium(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := parseUID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user_id"})
		return
	}

	ok, err := h.svc.CheckPremium(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"premium": ok,
	})
}
