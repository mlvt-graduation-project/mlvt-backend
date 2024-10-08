package handler

import (
	"mlvt/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MoMoPaymentHandler struct {
	momoPaymentService service.MoMoPaymentService
}

func NewMoMoPaymentHandler(momoPaymentService service.MoMoPaymentService) *MoMoPaymentHandler {
	return &MoMoPaymentHandler{momoPaymentService: momoPaymentService}
}

func (p *MoMoPaymentHandler) ProcessPayment(c *gin.Context) {
	var request struct {
		OrderID string `json:"order_id"`
		Amount  string `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Generate QR code for the payment
	qrCode, err := p.momoPaymentService.GeneratePaymentQRCode(request.OrderID, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Send back the QR code as an image
	c.Header("Content-Type", "image/png")
	c.Writer.Write(qrCode)
}

func (p *MoMoPaymentHandler) CheckPaymentStatus(c *gin.Context) {
	var request struct {
		OrderID string `json:"order_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	success, err := p.momoPaymentService.CheckPaymentStatus(request.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check payment status"})
		return
	}

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Payment not completed"})
	}
}

func (p *MoMoPaymentHandler) RefundPayment(c *gin.Context) {
	var request struct {
		OrderID string `json:"order_id"`
		Amount  string `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	refundResponse, err := p.momoPaymentService.RefundPayment(request.OrderID, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refund failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Refund successful", "data": refundResponse})
}
