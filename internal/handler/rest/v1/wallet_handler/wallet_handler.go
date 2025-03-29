package wallet_handler

import (
	"context"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/wallet_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	walletService wallet_service.WalletService
}

func NewWalletController(walletService wallet_service.WalletService) *WalletController {
	return &WalletController{
		walletService: walletService,
	}
}

// Deposit godoc
// @Summary      Deposit to the wallet
// @Description  Deposit a positive amount into a user's wallet
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Param        user_id  query     uint64 true  "User ID"
// @Param        amount   query     int    true   "Amount to deposit"
// @Success      200  {object}  response.MessageResponse  "Deposit successful"
// @Failure      400  {object}  response.ErrorResponse    "Invalid user ID or amount"
// @Failure      500  {object}  response.ErrorResponse    "Server error"
// @Router       /wallet/deposit [post]
func (wc *WalletController) Deposit(c *gin.Context) {
	userIDStr := c.Query("user_id")
	amountStr := c.Query("amount")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user_id"})
		return
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid amount"})
		return
	}

	err = wc.walletService.Deposit(context.Background(), userID, amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Deposit successful"})
}

// UseToken godoc
// @Summary      UseToken from the wallet
// @Description  UseToken a positive amount from a user's wallet if sufficient balance
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Param        user_id  query     uint64 true  "User ID"
// @Param        amount   query     int    true   "Amount to use token"
// @Success      200  {object}  response.MessageResponse  "UseTokenal successful"
// @Failure      400  {object}  response.ErrorResponse    "Invalid user ID, amount, or insufficient balance"
// @Failure      500  {object}  response.ErrorResponse    "Server error"
// @Router       /wallet/use-token [post]
func (wc *WalletController) UseToken(c *gin.Context) {
	userIDStr := c.Query("user_id")
	amountStr := c.Query("amount")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user_id"})
		return
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid amount"})
		return
	}

	err = wc.walletService.UseToken(context.Background(), userID, amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "UseTokenal successful"})
}

// GetBalance godoc
// @Summary      Get wallet balance
// @Description  Retrieve the current wallet balance for a user
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Param        user_id  query     uint64 true  "User ID"
// @Success      200  {object}  map[string]int64  "Example: {"balance": 100}"
// @Failure      400  {object}  response.ErrorResponse    "Invalid user_id"
// @Failure      500  {object}  response.ErrorResponse    "Server error"
// @Router       /wallet/balance [get]
func (wc *WalletController) GetBalance(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user_id"})
		return
	}

	balance, err := wc.walletService.GetBalance(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
