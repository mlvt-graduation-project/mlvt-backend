package service

import (
	"mlvt/internal/repo"

	"github.com/yeqown/go-qrcode/v2"
)

type PaymentService interface {
	GeneratePaymentQRCode(orderID, amount string) ([]byte, error)
	CheckPaymentStatus(orderID string) (bool, error)
	RefundPayment(orderID, amount string) (string, error)
}

type paymentService struct {
	momoRepo repo.MoMoRepo
}

func (p *paymentService) GeneratePaymentQRCode(orderID, amount string) ([]byte, error) {
	payURL, err := p.momoRepo.CreatePayment(orderID, amount)
	if err != nil {
		return nil, err
	}

	// Generate QR code from the payment URL
	png, err := qrcode.Encode(payURL, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return png, nil
}

func (p *paymentService) CheckPaymentStatus(orderID string) (bool, error) {
	return p.momoRepo.CheckPaymentStatus(orderID)
}

func (p *paymentService) RefundPayment(orderID, amount string) (string, error) {
	return p.momoRepo.RefundPayment(orderID, amount)
}
