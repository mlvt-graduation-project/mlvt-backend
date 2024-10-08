package service

import (
	"mlvt/internal/repo"

	qrcode "github.com/skip2/go-qrcode"
)

type MoMoPaymentService interface {
	GeneratePaymentQRCode(orderID, amount string) ([]byte, error)
	CheckPaymentStatus(orderID string) (bool, error)
	RefundPayment(orderID, amount string) (string, error)
}

type MoMopaymentService struct {
	momoRepo repo.MoMoRepo
}

func NewMoMoPaymentService(momoRepo repo.MoMoRepo) MoMoPaymentService {
	return &MoMopaymentService{momoRepo: momoRepo}
}

func (p *MoMopaymentService) GeneratePaymentQRCode(orderID, amount string) ([]byte, error) {
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

func (p *MoMopaymentService) CheckPaymentStatus(orderID string) (bool, error) {
	return p.momoRepo.CheckPaymentStatus(orderID)
}

func (p *MoMopaymentService) RefundPayment(orderID, amount string) (string, error) {
	return p.momoRepo.RefundPayment(orderID, amount)
}
