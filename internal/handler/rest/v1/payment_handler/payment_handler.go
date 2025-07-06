package payment_handler

import (
	"context"
	"mlvt/internal/entity"
	"mlvt/internal/pkg/response"
	"mlvt/internal/service/payment_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentController struct {
	paymentService payment_service.PaymentService
}

func NewPaymentController(paymentService payment_service.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// CreatePayment godoc
// @Summary      Create a new payment QR code
// @Description  Creates a new payment transaction and generates a QR code for the specified payment option
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        user_id  query     uint64 true  "User ID"
// @Param        option   query     string true  "Payment option (500k, 1m, 2m, 5m, 10m)"
// @Success      200  {object}  entity.PaymentTransaction  "Payment transaction created successfully"
// @Failure      400  {object}  response.ErrorResponse     "Invalid request parameters"
// @Failure      500  {object}  response.ErrorResponse     "Server error"
// @Router       /payment/create [post]
func (pc *PaymentController) CreatePayment(c *gin.Context) {
	userIDStr := c.Query("user_id")
	optionStr := c.Query("option")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user_id"})
		return
	}

	option := entity.PaymentOption(optionStr)
	if entity.GetPaymentOptionInfo(option) == nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid payment option"})
		return
	}

	payment, err := pc.paymentService.CreatePayment(context.Background(), userID, option)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetPayment godoc
// @Summary      Get payment by ID
// @Description  Retrieves a payment transaction by its ID
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        payment_id  path     string true  "Payment ID"
// @Success      200  {object}  entity.PaymentTransaction  "Payment transaction found"
// @Failure      400  {object}  response.ErrorResponse     "Invalid payment ID"
// @Failure      404  {object}  response.ErrorResponse     "Payment not found"
// @Failure      500  {object}  response.ErrorResponse     "Server error"
// @Router       /payment/{payment_id} [get]
func (pc *PaymentController) GetPayment(c *gin.Context) {
	paymentIDStr := c.Param("payment_id")
	paymentID, err := primitive.ObjectIDFromHex(paymentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid payment ID"})
		return
	}

	payment, err := pc.paymentService.GetPaymentByID(context.Background(), paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if payment == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetPaymentByTransactionID godoc
// @Summary      Get payment by transaction ID
// @Description  Retrieves a payment transaction by its transaction ID
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        transaction_id  path     string true  "Transaction ID"
// @Success      200  {object}  entity.PaymentTransaction  "Payment transaction found"
// @Failure      404  {object}  response.ErrorResponse     "Payment not found"
// @Failure      500  {object}  response.ErrorResponse     "Server error"
// @Router       /payment/transaction/{transaction_id} [get]
func (pc *PaymentController) GetPaymentByTransactionID(c *gin.Context) {
	transactionID := c.Param("transaction_id")

	payment, err := pc.paymentService.GetPaymentByTransactionID(context.Background(), transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if payment == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetUserPayments godoc
// @Summary      Get user's payment history
// @Description  Retrieves all payment transactions for a specific user
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        user_id  query     uint64 true  "User ID"
// @Success      200  {array}   entity.PaymentTransaction  "User payments"
// @Failure      400  {object}  response.ErrorResponse     "Invalid user ID"
// @Failure      500  {object}  response.ErrorResponse     "Server error"
// @Router       /payment/user-payments [get]
func (pc *PaymentController) GetUserPayments(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user_id"})
		return
	}

	payments, err := pc.paymentService.GetUserPayments(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// ConfirmPayment godoc
// @Summary      Confirm a payment transaction
// @Description  Marks a payment as completed and adds tokens to user's wallet
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        transaction_id  path     string true  "Transaction ID"
// @Success      200  {object}  response.MessageResponse  "Payment confirmed successfully"
// @Failure      400  {object}  response.ErrorResponse    "Invalid transaction ID or payment not pending"
// @Failure      404  {object}  response.ErrorResponse    "Payment not found"
// @Failure      500  {object}  response.ErrorResponse    "Server error"
// @Router       /payment/confirm/{transaction_id} [post]
func (pc *PaymentController) ConfirmPayment(c *gin.Context) {
	transactionID := c.Param("transaction_id")

	err := pc.paymentService.ConfirmPayment(context.Background(), transactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Payment confirmed successfully"})
}

// CancelPayment godoc
// @Summary      Cancel a payment transaction
// @Description  Marks a payment as cancelled
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        payment_id  path     string true  "Payment ID"
// @Success      200  {object}  response.MessageResponse  "Payment cancelled successfully"
// @Failure      400  {object}  response.ErrorResponse    "Invalid payment ID"
// @Failure      500  {object}  response.ErrorResponse    "Server error"
// @Router       /payment/cancel/{payment_id} [post]
func (pc *PaymentController) CancelPayment(c *gin.Context) {
	paymentIDStr := c.Param("payment_id")
	paymentID, err := primitive.ObjectIDFromHex(paymentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid payment ID"})
		return
	}

	err = pc.paymentService.CancelPayment(context.Background(), paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Payment cancelled successfully"})
}

// GetPaymentOptions godoc
// @Summary      Get available payment options
// @Description  Retrieves all available payment options with their token amounts and VND prices
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Success      200  {array}   entity.PaymentOptionInfo  "Payment options"
// @Router       /payment/options [get]
func (pc *PaymentController) GetPaymentOptions(c *gin.Context) {
	options := pc.paymentService.GetPaymentOptions()
	c.JSON(http.StatusOK, options)
}

// GetPendingPayments godoc
// @Summary      Get pending payments
// @Description  Retrieves all pending payment transactions (Admin only)
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Success      200  {array}   entity.PaymentTransaction  "Pending payments"
// @Failure      500  {object}  response.ErrorResponse     "Server error"
// @Router       /payment/pending [get]
func (pc *PaymentController) GetPendingPayments(c *gin.Context) {
	payments, err := pc.paymentService.GetPendingPayments(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}
